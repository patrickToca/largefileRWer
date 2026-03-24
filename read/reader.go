package read

import (
	"bufio"
	"context"
	"crypto/sha256"
	"crypto/sha512"
	"fmt"
	"io"
	"log/slog"
	"os"
	"sync"

	"pmjtoca/largefileRWer/config"
	"pmjtoca/largefileRWer/internal/checkpoint"
	"pmjtoca/largefileRWer/internal/metrics"
)

// Chunk represents a piece of data read from the file
type Chunk struct {
	ID       int
	Offset   int64
	Data     []byte
	Checksum []byte
	Error    error
}

// Reader handles efficient reading of large files
type Reader struct {
	cfg        *config.Config
	metrics    *metrics.Metrics
	checkpoint *checkpoint.Manager
	bufferPool *sync.Pool
	workerPool chan struct{}
	ctx        context.Context
	cancel     context.CancelFunc
	hashFunc   func([]byte) []byte
}

// NewReader creates a new file reader
func NewReader(cfg *config.Config) (*Reader, error) {
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	// Create buffer pool with aligned buffers for direct I/O
	bufferPool := &sync.Pool{
		New: func() interface{} {
			buf := make([]byte, cfg.GetBufferSizeBytes()+cfg.DirectIOAlignment)
			aligned := buf[:cfg.GetBufferSizeBytes()]
			return &aligned
		},
	}

	// Setup hash function based on configuration
	var hashFunc func([]byte) []byte
	switch cfg.ChecksumAlgorithm {
	case "sha256":
		hashFunc = func(data []byte) []byte {
			h := sha256.Sum256(data)
			return h[:]
		}
	case "sha512":
		hashFunc = func(data []byte) []byte {
			h := sha512.Sum512(data)
			return h[:]
		}
	default:
		hashFunc = func(data []byte) []byte {
			h := sha256.Sum256(data)
			return h[:]
		}
	}

	// Initialize checkpoint manager
	checkpointMgr := checkpoint.NewManager(cfg.GetCheckpointPath())
	if err := checkpointMgr.Load(); err != nil {
		slog.Warn("Failed to load checkpoint", "error", err)
	}

	return &Reader{
		cfg:        cfg,
		metrics:    metrics.NewMetrics(),
		checkpoint: checkpointMgr,
		bufferPool: bufferPool,
		workerPool: make(chan struct{}, cfg.Workers),
		ctx:        ctx,
		cancel:     cancel,
		hashFunc:   hashFunc,
	}, nil
}

// ReadAll reads the entire file and returns a channel of chunks
func (r *Reader) ReadAll(ctx context.Context) (<-chan *Chunk, <-chan error) {
	chunkChan := make(chan *Chunk, r.cfg.QueueSize)
	errChan := make(chan error, 1)

	go func() {
		defer close(chunkChan)
		defer close(errChan)

		fileInfo, err := os.Stat(r.cfg.InputPath)
		if err != nil {
			errChan <- fmt.Errorf("failed to stat file: %w", err)
			return
		}
		fileSize := fileInfo.Size()

		slog.Info("Starting read",
			"path", r.cfg.InputPath,
			"size_gb", float64(fileSize)/(1024*1024*1024),
			"workers", r.cfg.Workers)

		// Choose strategy based on file size and configuration
		if fileSize < 1024*1024*1024 {
			r.readSmallFile(ctx, chunkChan, errChan)
		} else if r.cfg.UseMmap {
			r.readWithMmap(ctx, fileSize, chunkChan, errChan)
		} else {
			r.readWithParallelChunks(ctx, fileSize, chunkChan, errChan)
		}
	}()

	return chunkChan, errChan
}

// readSmallFile handles files under 1GB with simple buffered I/O
func (r *Reader) readSmallFile(ctx context.Context, chunkChan chan<- *Chunk, errChan chan<- error) {
	file, err := os.Open(r.cfg.InputPath)
	if err != nil {
		errChan <- err
		return
	}
	defer file.Close()

	if r.cfg.UseDirectIO {
		r.enableDirectIO(file)
	}

	reader := bufio.NewReaderSize(file, r.cfg.GetBufferSizeBytes())
	bufPtr := r.bufferPool.Get().(*[]byte)
	defer r.bufferPool.Put(bufPtr)

	buffer := *bufPtr
	var offset int64
	chunkID := 0

	for {
		select {
		case <-ctx.Done():
			errChan <- ctx.Err()
			return
		default:
		}

		n, err := reader.Read(buffer)
		if n > 0 {
			data := make([]byte, n)
			copy(data, buffer[:n])

			var checksum []byte
			if r.cfg.VerifyChecksum {
				checksum = r.hashFunc(data)
			}

			select {
			case chunkChan <- &Chunk{
				ID:       chunkID,
				Offset:   offset,
				Data:     data,
				Checksum: checksum,
			}:
			case <-ctx.Done():
				return
			}

			offset += int64(n)
			chunkID++
			r.metrics.AddBytesRead(int64(n))
			r.metrics.AddChunkRead()
		}

		if err == io.EOF {
			break
		}
		if err != nil {
			errChan <- err
			return
		}
	}
}

// readWithParallelChunks reads file by splitting into parallel chunks
func (r *Reader) readWithParallelChunks(ctx context.Context, fileSize int64,
	chunkChan chan<- *Chunk, errChan chan<- error) {

	file, err := os.Open(r.cfg.InputPath)
	if err != nil {
		errChan <- err
		return
	}
	defer file.Close()

	if r.cfg.UseDirectIO {
		r.enableDirectIO(file)
	}

	chunkSize := r.cfg.GetChunkSizeBytes()
	numChunks := (fileSize + chunkSize - 1) / chunkSize

	slog.Debug("Splitting into chunks",
		"num_chunks", numChunks,
		"chunk_size_mb", r.cfg.ChunkSizeMB)

	var wg sync.WaitGroup
	var once sync.Once
	var firstErr error
	var errMutex sync.Mutex

	for chunkID := 0; chunkID < int(numChunks); chunkID++ {
		select {
		case <-ctx.Done():
			errChan <- ctx.Err()
			return
		default:
		}

		// Skip already processed chunks
		if r.checkpoint.IsCompleted(chunkID) {
			slog.Debug("Skipping completed chunk", "chunk_id", chunkID)
			continue
		}

		offset := int64(chunkID) * chunkSize
		length := chunkSize
		if offset+length > fileSize {
			length = fileSize - offset
		}

		// Acquire worker slot
		select {
		case r.workerPool <- struct{}{}:
		case <-ctx.Done():
			return
		}

		wg.Add(1)
		go func(id int, off, ln int64) {
			defer wg.Done()
			defer func() { <-r.workerPool }()

			chunk, err := r.readChunk(file, id, off, ln)
			if err != nil {
				errMutex.Lock()
				if firstErr == nil {
					firstErr = err
					once.Do(func() {
						errChan <- err
						r.cancel()
					})
				}
				errMutex.Unlock()
				return
			}

			select {
			case chunkChan <- chunk:
			case <-ctx.Done():
				return
			}

			r.metrics.AddChunkRead()
			r.checkpoint.MarkCompleted(id)

			// Periodic checkpoint save
			if id%r.cfg.CheckpointInterval == 0 {
				if err := r.checkpoint.Save(); err != nil {
					slog.Warn("Failed to save checkpoint", "error", err)
				}
			}
		}(chunkID, offset, length)
	}

	wg.Wait()

	if firstErr != nil {
		errChan <- firstErr
	}
}

// readChunk reads a single chunk from the file
func (r *Reader) readChunk(file *os.File, chunkID int, offset, length int64) (*Chunk, error) {
	bufPtr := r.bufferPool.Get().(*[]byte)
	defer r.bufferPool.Put(bufPtr)

	buffer := *bufPtr
	if int64(len(buffer)) < length {
		buffer = make([]byte, length)
	}

	data := buffer[:length]
	n, err := file.ReadAt(data, offset)
	if err != nil && err != io.EOF {
		return nil, fmt.Errorf("chunk %d read error: %w", chunkID, err)
	}

	// Create a copy to avoid buffer reuse issues
	dataCopy := make([]byte, n)
	copy(dataCopy, data[:n])

	var checksum []byte
	if r.cfg.VerifyChecksum {
		checksum = r.hashFunc(dataCopy)
	}

	r.metrics.AddBytesRead(int64(n))

	return &Chunk{
		ID:       chunkID,
		Offset:   offset,
		Data:     dataCopy,
		Checksum: checksum,
	}, nil
}

// readWithMmap uses memory-mapped I/O for random access patterns
func (r *Reader) readWithMmap(ctx context.Context, fileSize int64,
	chunkChan chan<- *Chunk, errChan chan<- error) {

	slog.Info("Using mmap for reading", "file_size_gb", float64(fileSize)/(1024*1024*1024))

	// Note: mmap implementation would go here
	// For now, fall back to parallel chunks
	slog.Warn("mmap not fully implemented, falling back to parallel chunks")
	r.readWithParallelChunks(ctx, fileSize, chunkChan, errChan)
}

// GetMetrics returns the metrics collector
func (r *Reader) GetMetrics() *metrics.Metrics {
	return r.metrics
}

// Close cleans up resources
func (r *Reader) Close() error {
	r.cancel()
	if err := r.checkpoint.Save(); err != nil {
		slog.Warn("Failed to save final checkpoint", "error", err)
	}
	return nil
}
