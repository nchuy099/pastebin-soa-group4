package consumer

import (
	"encoding/json"
	"log"
	"sync"

	"get-paste-worker/model"
	"get-paste-worker/repository"

	"github.com/rabbitmq/amqp091-go"
)

// Worker represents a worker in the worker pool
type Worker struct {
	ID            int
	WorkerPool    chan chan amqp091.Delivery
	JobChannel    chan amqp091.Delivery
	Quit          chan bool
	ProcessedJobs int
}

// WorkerPool represents a pool of workers
type WorkerPool struct {
	NumWorkers    int
	WorkerQueue   chan chan amqp091.Delivery
	Workers       []*Worker
	MessageSource <-chan amqp091.Delivery
	Quit          chan bool
	wg            sync.WaitGroup
}

// NewWorker creates a new worker
func NewWorker(id int, workerPool chan chan amqp091.Delivery) *Worker {
	return &Worker{
		ID:         id,
		WorkerPool: workerPool,
		JobChannel: make(chan amqp091.Delivery),
		Quit:       make(chan bool),
	}
}

// Start starts the worker
func (w *Worker) Start() {
	go func() {
		for {
			// Register the current worker into the worker pool
			w.WorkerPool <- w.JobChannel

			select {
			case job := <-w.JobChannel:
				// Received a job, process it
				log.Printf("Worker %d processing message", w.ID)
				if err := processMessage(job); err != nil {
					log.Printf("Worker %d error processing message: %v", w.ID, err)
					job.Nack(false, true) // Reject and requeue the message
				} else {
					job.Ack(false) // Acknowledge the message
					w.ProcessedJobs++
					log.Printf("Worker %d processed message (total: %d)", w.ID, w.ProcessedJobs)
				}

			case <-w.Quit:
				// We were told to stop
				log.Printf("Worker %d stopping", w.ID)
				return
			}
		}
	}()
}

// Stop stops the worker
func (w *Worker) Stop() {
	go func() {
		w.Quit <- true
	}()
}

// NewWorkerPool creates a new worker pool
func NewWorkerPool(numWorkers int, messageSource <-chan amqp091.Delivery) *WorkerPool {
	workerPool := &WorkerPool{
		NumWorkers:    numWorkers,
		WorkerQueue:   make(chan chan amqp091.Delivery, numWorkers),
		Workers:       make([]*Worker, numWorkers),
		MessageSource: messageSource,
		Quit:          make(chan bool),
	}

	// Create and start workers
	for i := 0; i < numWorkers; i++ {
		worker := NewWorker(i, workerPool.WorkerQueue)
		workerPool.Workers[i] = worker
		worker.Start()
	}

	return workerPool
}

// Start starts the worker pool
func (wp *WorkerPool) Start() {
	wp.wg.Add(1)
	go func() {
		defer wp.wg.Done()
		for {
			select {
			case message := <-wp.MessageSource:
				// Received a message, dispatch to a worker
				go func(message amqp091.Delivery) {
					// Get a worker from the pool
					jobChannel := <-wp.WorkerQueue

					// Send the job to the worker
					jobChannel <- message
				}(message)

			case <-wp.Quit:
				// We were told to stop
				log.Println("Worker pool stopping")
				return
			}
		}
	}()
}

// Stop stops the worker pool
func (wp *WorkerPool) Stop() {
	// Stop all workers
	for _, worker := range wp.Workers {
		worker.Stop()
	}

	// Stop the dispatcher
	wp.Quit <- true
	wp.wg.Wait()
}

// processMessage processes a message from the queue
func processMessage(delivery amqp091.Delivery) error {
	// Parse the message
	var message ViewUpdateMessage
	err := json.Unmarshal(delivery.Body, &message)
	if err != nil {
		log.Printf("Error unmarshaling message: %v", err)
		return err
	}

	// Create paste view object
	pasteView := &model.PasteViews{
		PasteID:  message.PasteID,
		ViewedAt: message.ViewedAt,
	}

	// Update the database using batch insert
	err = repository.AddPasteViewBatch(pasteView)
	if err != nil {
		log.Printf("Error adding paste view: %v", err)
		return err
	}

	log.Printf("Successfully processed view update for paste %s", message.PasteID)
	return nil
}
