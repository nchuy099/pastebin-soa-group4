package repository

import (
	"log"
	"strings"
	"sync"
	"time"

	"create-paste-worker/db"
	"create-paste-worker/model"
)

// Configuration variables for batch processing
var (
	BatchSize        = 50 // Default number of pastes per batch
	BatchTimeoutSecs = 5  // Default timeout in seconds
)

// BatchManager handles batch inserts of pastes
type BatchManager struct {
	mu          sync.Mutex
	pastes      []*model.Paste
	lastFlushAt time.Time
	batchTimer  *time.Timer
	batchSize   int
}

var batchManager *BatchManager
var once sync.Once

// GetBatchManager returns the singleton batch manager instance
func GetBatchManager() *BatchManager {
	once.Do(func() {
		batchManager = &BatchManager{
			pastes:      make([]*model.Paste, 0, BatchSize),
			lastFlushAt: time.Now(),
			batchSize:   BatchSize,
		}
		// Start the timer for time-based flushing
		batchManager.batchTimer = time.AfterFunc(time.Duration(BatchTimeoutSecs)*time.Second, func() {
			batchManager.FlushBatch()
			batchManager.resetTimer()
		})
	})
	return batchManager
}

// resetTimer resets the batch timer
func (bm *BatchManager) resetTimer() {
	bm.batchTimer.Reset(time.Duration(BatchTimeoutSecs) * time.Second)
}

// AddPaste adds a paste to the batch and flushes if needed
func (bm *BatchManager) AddPaste(paste *model.Paste) error {
	bm.mu.Lock()
	defer bm.mu.Unlock()

	bm.pastes = append(bm.pastes, paste)

	// If we've reached the batch size, flush immediately
	if len(bm.pastes) >= bm.batchSize {
		return bm.flushBatchLocked()
	}

	return nil
}

// FlushBatch forces flushing of the current batch
func (bm *BatchManager) FlushBatch() error {
	bm.mu.Lock()
	defer bm.mu.Unlock()

	return bm.flushBatchLocked()
}

// flushBatchLocked performs the actual batch insert (assumes lock is held)
func (bm *BatchManager) flushBatchLocked() error {
	if len(bm.pastes) == 0 {
		return nil // Nothing to flush
	}

	log.Printf("Flushing batch of %d pastes", len(bm.pastes))

	// Build batch insert query
	query := `
		INSERT INTO paste (id, content, title, language, created_at, expires_at, visibility)
		VALUES 
	`

	// Create value placeholders and flatten values
	placeholders := make([]string, 0, len(bm.pastes))
	args := make([]interface{}, 0, len(bm.pastes)*7)

	for _, paste := range bm.pastes {
		placeholders = append(placeholders, "(?, ?, ?, ?, ?, ?, ?)")
		args = append(args,
			paste.ID,
			paste.Content,
			paste.Title,
			paste.Language,
			paste.CreatedAt,
			paste.ExpiresAt,
			paste.Visibility,
		)
	}

	// Complete the query with placeholders
	query += strings.Join(placeholders, ", ")

	// Execute the batch insert
	_, err := db.DB.Exec(query, args...)
	if err != nil {
		log.Printf("Error batch saving pastes: %v", err.Error())
		return err
	}

	log.Printf("Successfully batch inserted %d pastes", len(bm.pastes))

	// Clear the batch
	bm.pastes = bm.pastes[:0]
	bm.lastFlushAt = time.Now()

	return nil
}

// SavePaste saves a new paste to the database (single insert)
func SavePaste(paste *model.Paste) error {
	query := `
		INSERT INTO paste (id, content, title, language, created_at, expires_at, visibility)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`

	_, err := db.DB.Exec(
		query,
		paste.ID,
		paste.Content,
		paste.Title,
		paste.Language,
		paste.CreatedAt,
		paste.ExpiresAt,
		paste.Visibility,
	)

	if err != nil {
		log.Printf("Error saving paste: %v", err.Error())
		return err
	}

	return nil
}

// SavePasteBatch adds a paste to the batch insert queue
func SavePasteBatch(paste *model.Paste) error {
	return GetBatchManager().AddPaste(paste)
}

// Shutdown ensures any pending pastes are flushed before shutdown
func Shutdown() error {
	if batchManager != nil {
		return batchManager.FlushBatch()
	}
	return nil
}
