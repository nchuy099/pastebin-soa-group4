package repository

import (
	"log"
	"strings"
	"sync"
	"time"

	"get-paste-worker/db"
	"get-paste-worker/model"
)

// Configuration variables for batch processing
var (
	BatchSize        = 50 // Default number of paste views per batch
	BatchTimeoutSecs = 5  // Default timeout in seconds
)

// BatchManager handles batch inserts of paste views
type BatchManager struct {
	mu          sync.Mutex
	pasteViews  []*model.PasteViews
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
			pasteViews:  make([]*model.PasteViews, 0, BatchSize),
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

// AddPasteView adds a paste view to the batch and flushes if needed
func (bm *BatchManager) AddPasteView(pasteView *model.PasteViews) error {
	bm.mu.Lock()
	defer bm.mu.Unlock()

	bm.pasteViews = append(bm.pasteViews, pasteView)

	// If we've reached the batch size, flush immediately
	if len(bm.pasteViews) >= bm.batchSize {
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
	if len(bm.pasteViews) == 0 {
		return nil // Nothing to flush
	}

	log.Printf("Flushing batch of %d paste views", len(bm.pasteViews))

	// Build batch insert query
	query := `
		INSERT INTO paste_views (paste_id, viewed_at)
		VALUES 
	`

	// Create value placeholders and flatten values
	placeholders := make([]string, 0, len(bm.pasteViews))
	args := make([]interface{}, 0, len(bm.pasteViews)*2)

	for _, view := range bm.pasteViews {
		placeholders = append(placeholders, "(?, ?)")
		args = append(args,
			view.PasteID,
			view.ViewedAt,
		)
	}

	// Complete the query with placeholders
	query += strings.Join(placeholders, ", ")

	// Execute the batch insert
	_, err := db.DB.Exec(query, args...)
	if err != nil {
		log.Printf("Error batch saving paste views: %v", err.Error())
		return err
	}

	log.Printf("Successfully batch inserted %d paste views", len(bm.pasteViews))

	// Clear the batch
	bm.pasteViews = bm.pasteViews[:0]
	bm.lastFlushAt = time.Now()

	return nil
}

// AddPasteViewDirect adds a paste view record directly to the database
func AddPasteViewDirect(pasteView *model.PasteViews) error {
	query := `
		INSERT INTO paste_views (paste_id, viewed_at)
		VALUES (?, ?)
	`
	_, err := db.DB.Exec(query, pasteView.PasteID, pasteView.ViewedAt)
	if err != nil {
		log.Printf("Error adding paste view: %v", err.Error())
	}
	return nil
}

// AddPasteViewBatch adds a paste view to the batch insert queue
func AddPasteViewBatch(pasteView *model.PasteViews) error {
	return GetBatchManager().AddPasteView(pasteView)
}

// Shutdown ensures any pending paste views are flushed before shutdown
func Shutdown() error {
	if batchManager != nil {
		return batchManager.FlushBatch()
	}
	return nil
}
