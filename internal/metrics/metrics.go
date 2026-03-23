package metrics

import (
	"sync/atomic"
	"time"
)

// Metrics tracks performance metrics
type Metrics struct {
	BytesRead     atomic.Int64
	BytesWritten  atomic.Int64
	ChunksRead    atomic.Int64
	ChunksWritten atomic.Int64
	StartTime     time.Time
}

// NewMetrics creates a new metrics collector
func NewMetrics() *Metrics {
	return &Metrics{
		StartTime: time.Now(),
	}
}

// AddBytesRead adds bytes to read counter
func (m *Metrics) AddBytesRead(n int64) {
	m.BytesRead.Add(n)
}

// AddBytesWritten adds bytes to write counter
func (m *Metrics) AddBytesWritten(n int64) {
	m.BytesWritten.Add(n)
}

// AddChunkRead increments chunk counter
func (m *Metrics) AddChunkRead() {
	m.ChunksRead.Add(1)
}

// AddChunkWritten increments chunk counter
func (m *Metrics) AddChunkWritten() {
	m.ChunksWritten.Add(1)
}

// GetBytesRead returns total bytes read
func (m *Metrics) GetBytesRead() int64 {
	return m.BytesRead.Load()
}

// GetBytesWritten returns total bytes written
func (m *Metrics) GetBytesWritten() int64 {
	return m.BytesWritten.Load()
}

// GetSpeedMBps returns current speed in MB/s
func (m *Metrics) GetSpeedMBps(read bool) float64 {
	elapsed := time.Since(m.StartTime)
	if elapsed == 0 {
		return 0
	}
	var bytes int64
	if read {
		bytes = m.BytesRead.Load()
	} else {
		bytes = m.BytesWritten.Load()
	}
	return float64(bytes) / elapsed.Seconds() / (1024 * 1024)
}

// GetElapsed returns elapsed time
func (m *Metrics) GetElapsed() time.Duration {
	return time.Since(m.StartTime)
}
