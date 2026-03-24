package write

import (
	"context"
	"crypto/sha256"
	"crypto/sha512"
	"fmt"
	"log/slog"
	"os"
	"sync"

	"pmjtoca/largefileRWer/config"
	"pmjtoca/largefileRWer/internal/metrics"
)

// WriteChunk represents a chunk to be written
type WriteChunk struct {
	ID       int
	Offset   int64
	Data     []byte
	Checksum []byte
}

// Writer handles efficient writing of large files
type Writer struct {
	cfg           *config.Config
	metrics       *metrics.Metrics
	bufferPool    *sync.Pool
	workerPool    chan struct{}
	pendingChunks map[int]*WriteChunk
	nextChunkID   int
	mu            sync.RWMutex
	ctx           context.Context
	cancel        context.CancelFunc
	hashFunc      func([]byte) []byte
}

// NewWriter creates a new file writer
func NewWriter(cfg *config.Config) (*Writer, error) {
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	bufferPool := &sync.Pool{
		New: func() interface{} {
			buf := make([]byte, cfg.GetBufferSizeBytes()+cfg.DirectIOAlignment)
			aligned := buf[:cfg.GetBufferSizeBytes()]
			return &aligned
		},
	}

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

	return &Writer{
		cfg:           cfg,
		metrics:       metrics.NewMetrics(),
		bufferPool:    bufferPool,
		workerPool:    make(chan struct{}, cfg.Workers),
		pendingChunks: make(map[int]*WriteChunk),
		ctx:           ctx,
		cancel:        cancel,
		hashFunc:      hashFunc,
	}, nil
}

// WriteChunksOrdered writes chunks in the correct order, handling out-of-order arrivals
func (w *Writer) WriteChunksOrdered(ctx context.Context, chunkChan <-chan *WriteChunk) error {
	file, err := os.OpenFile(w.cfg.OutputPath, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer file.Close()

	if w.cfg.UseDirectIO {
		w.enableDirectIO(file)
	}

	slog.Info("Starting ordered write", "path", w.cfg.OutputPath)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case chunk, ok := <-chunkChan:
			if !ok {
				return nil
			}
			if err := w.processChunk(file, chunk); err != nil {
				return err
			}
		}
	}
}

// processChunk handles incoming chunks, buffering out-of-order chunks
func (w *Writer) processChunk(file *os.File, chunk *WriteChunk) error {
	w.mu.Lock()
	defer w.mu.Unlock()

	// Verify checksum if enabled
	if w.cfg.VerifyChecksum {
		calculated := w.hashFunc(chunk.Data)
		if len(calculated) != len(chunk.Checksum) {
			return fmt.Errorf("checksum length mismatch for chunk %d", chunk.ID)
		}
		for i := range calculated {
			if calculated[i] != chunk.Checksum[i] {
				return fmt.Errorf("checksum mismatch for chunk %d", chunk.ID)
			}
		}
	}

	// Store chunk
	w.pendingChunks[chunk.ID] = chunk

	// Write all consecutive chunks
	for {
		nextChunk, exists := w.pendingChunks[w.nextChunkID]
		if !exists {
			break
		}

		_, err := file.WriteAt(nextChunk.Data, nextChunk.Offset)
		if err != nil {
			return fmt.Errorf("failed to write chunk %d: %w", nextChunk.ID, err)
		}

		w.metrics.AddBytesWritten(int64(len(nextChunk.Data)))
		w.metrics.AddChunkWritten()

		// Log progress periodically
		if w.nextChunkID%w.cfg.CheckpointInterval == 0 {
			slog.Debug("Write progress",
				"chunks_written", w.nextChunkID,
				"bytes_written_mb", w.metrics.GetBytesWritten()/(1024*1024))
		}

		delete(w.pendingChunks, w.nextChunkID)
		w.nextChunkID++
	}

	return nil
}

// WriteChunk writes a single chunk (for non-ordered writes)
func (w *Writer) WriteChunk(ctx context.Context, chunk *WriteChunk) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case w.workerPool <- struct{}{}:
		defer func() { <-w.workerPool }()
	}

	// Verify checksum if enabled
	if w.cfg.VerifyChecksum {
		calculated := w.hashFunc(chunk.Data)
		if len(calculated) != len(chunk.Checksum) {
			return fmt.Errorf("checksum length mismatch for chunk %d", chunk.ID)
		}
		for i := range calculated {
			if calculated[i] != chunk.Checksum[i] {
				return fmt.Errorf("checksum mismatch for chunk %d", chunk.ID)
			}
		}
	}

	file, err := os.OpenFile(w.cfg.OutputPath, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	if w.cfg.UseDirectIO {
		w.enableDirectIO(file)
	}

	_, err = file.WriteAt(chunk.Data, chunk.Offset)
	if err != nil {
		return err
	}

	w.metrics.AddBytesWritten(int64(len(chunk.Data)))
	w.metrics.AddChunkWritten()

	return nil
}

// WriteChunksParallel writes chunks in parallel (for independent chunks)
func (w *Writer) WriteChunksParallel(ctx context.Context, chunks []*WriteChunk) error {
	if len(chunks) == 0 {
		return nil
	}

	var wg sync.WaitGroup
	errChan := make(chan error, len(chunks))

	// Limit concurrent writes
	semaphore := make(chan struct{}, w.cfg.Workers)

	for _, chunk := range chunks {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		wg.Add(1)
		go func(c *WriteChunk) {
			defer wg.Done()

			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			if err := w.WriteChunk(ctx, c); err != nil {
				errChan <- err
			}
		}(chunk)
	}

	wg.Wait()
	close(errChan)

	// Check for errors
	for err := range errChan {
		if err != nil {
			return err
		}
	}

	return nil
}

// WriteSequential writes data sequentially (simplest approach)
func (w *Writer) WriteSequential(ctx context.Context, data []byte) error {
	file, err := os.OpenFile(w.cfg.OutputPath, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	if w.cfg.UseDirectIO {
		w.enableDirectIO(file)
	}

	bufPtr := w.bufferPool.Get().(*[]byte)
	defer w.bufferPool.Put(bufPtr)

	buffer := *bufPtr
	var written int64

	for written < int64(len(data)) {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		toWrite := int64(len(buffer))
		if written+toWrite > int64(len(data)) {
			toWrite = int64(len(data)) - written
		}

		copy(buffer, data[written:written+toWrite])

		n, err := file.WriteAt(buffer[:toWrite], written)
		if err != nil {
			return err
		}

		written += int64(n)
		w.metrics.AddBytesWritten(int64(n))
	}

	return nil
}

// CreateEmptyFile creates an empty file with the given size (sparse file)
func (w *Writer) CreateEmptyFile(ctx context.Context, sizeBytes int64) error {
	file, err := os.OpenFile(w.cfg.OutputPath, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	// Create sparse file by seeking to end-1 and writing a single byte
	_, err = file.Seek(sizeBytes-1, 0)
	if err != nil {
		return err
	}

	_, err = file.Write([]byte{0})
	return err
}

// Flush ensures all data is written to disk
func (w *Writer) Flush() error {
	file, err := os.OpenFile(w.cfg.OutputPath, os.O_RDONLY, 0)
	if err != nil {
		return err
	}
	defer file.Close()
	return file.Sync()
}

// GetMetrics returns the metrics collector
func (w *Writer) GetMetrics() *metrics.Metrics {
	return w.metrics
}

// Close cleans up resources
func (w *Writer) Close() error {
	w.cancel()
	return nil
}