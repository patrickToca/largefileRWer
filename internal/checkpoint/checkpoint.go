package checkpoint

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"sync"
)

// Manager handles checkpoint operations
type Manager struct {
	path           string
	mu             sync.RWMutex
	completed      map[int]bool
	lastCheckpoint int
}

// NewManager creates a new checkpoint manager
func NewManager(path string) *Manager {
	return &Manager{
		path:      path,
		completed: make(map[int]bool),
	}
}

// Load loads checkpoint from disk
func (cm *Manager) Load() error {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	data, err := os.ReadFile(cm.path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("failed to read checkpoint: %w", err)
	}

	if err := json.Unmarshal(data, &cm.completed); err != nil {
		return fmt.Errorf("failed to parse checkpoint: %w", err)
	}

	slog.Info("Loaded checkpoint", "completed_chunks", len(cm.completed))
	return nil
}

// Save saves checkpoint to disk
func (cm *Manager) Save() error {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	data, err := json.Marshal(cm.completed)
	if err != nil {
		return fmt.Errorf("failed to marshal checkpoint: %w", err)
	}

	tmpFile := cm.path + ".tmp"
	if err := os.WriteFile(tmpFile, data, 0644); err != nil {
		return fmt.Errorf("failed to write checkpoint: %w", err)
	}

	if err := os.Rename(tmpFile, cm.path); err != nil {
		return fmt.Errorf("failed to rename checkpoint: %w", err)
	}

	return nil
}

// MarkCompleted marks a chunk as completed
func (cm *Manager) MarkCompleted(chunkID int) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	cm.completed[chunkID] = true
}

// IsCompleted checks if a chunk is completed
func (cm *Manager) IsCompleted(chunkID int) bool {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	return cm.completed[chunkID]
}
