package main

import (
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
	"time"

	"get-paste-worker/config"
	"get-paste-worker/db"
	"get-paste-worker/model"
	"get-paste-worker/repository"

	"github.com/rabbitmq/amqp091-go"
)

// Worker configuration
const (
	DefaultNumWorkers = 10 // Default number of worker goroutines
	QueueName         = "view_updates_queue"
	ExchangeName      = "view_updates"
	RoutingKey        = "view.update"
)

// ViewUpdateMessage structure for RabbitMQ messages
type ViewUpdateMessage struct {
	PasteID    string    `json:"paste_id"`
	ViewedAt   time.Time `json:"viewed_at"`
	RemoteAddr string    `json:"remote_addr,omitempty"`
	UserAgent  string    `json:"user_agent,omitempty"`
}

func main() {
	log.Println("Starting get paste worker...")

	// Load environment variables and initialize database
	config.LoadEnv()
	db.InitDB()

	// Get number of workers from environment variable
	numWorkers := DefaultNumWorkers
	if numWorkersEnv := os.Getenv("NUM_WORKERS"); numWorkersEnv != "" {
		if n, err := strconv.Atoi(numWorkersEnv); err == nil && n > 0 {
			numWorkers = n
		} else if err != nil {
			log.Printf("Invalid NUM_WORKERS value: %s, using default: %d", numWorkersEnv, DefaultNumWorkers)
		}
	}
	log.Printf("Configured to use %d worker goroutines", numWorkers)

	// Set up RabbitMQ connection
	rabbitURL := os.Getenv("RABBITMQ_URL")
	if rabbitURL == "" {
		rabbitURL = "amqp://guest:guest@localhost:5672/"
	}

	conn, err := amqp091.Dial(rabbitURL)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %v", err)
	}
	defer ch.Close()

	// Declare exchange
	err = ch.ExchangeDeclare(
		ExchangeName, // name
		"direct",     // type
		true,         // durable
		false,        // auto-deleted
		false,        // internal
		false,        // no-wait
		nil,          // arguments
	)
	if err != nil {
		log.Fatalf("Failed to declare an exchange: %v", err)
	}

	// Declare queue
	q, err := ch.QueueDeclare(
		QueueName, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		log.Fatalf("Failed to declare a queue: %v", err)
	}

	// Bind queue to exchange
	err = ch.QueueBind(
		q.Name,       // queue name
		RoutingKey,   // routing key
		ExchangeName, // exchange
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to bind queue: %v", err)
	}

	// Set up consumer
	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		log.Fatalf("Failed to register a consumer: %v", err)
	}

	// Set up worker pool
	var wg sync.WaitGroup
	jobs := make(chan amqp091.Delivery)

	// Start worker goroutines
	for w := 1; w <= numWorkers; w++ {
		wg.Add(1)
		go worker(w, jobs, &wg)
	}

	// Set up signal handling for graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Start dispatcher
	go func() {
		for msg := range msgs {
			jobs <- msg
		}
	}()

	log.Printf("Get paste worker is running with %d goroutines", numWorkers)
	log.Printf("Press CTRL+C to exit")

	// Wait for shutdown signal
	<-quit
	log.Println("Shutting down worker...")

	// Close channel to stop workers
	close(jobs)

	// Wait for all workers to finish
	wg.Wait()
	log.Println("Worker shutdown complete")
}

func worker(id int, jobs <-chan amqp091.Delivery, wg *sync.WaitGroup) {
	defer wg.Done()

	log.Printf("Worker %d starting", id)

	for msg := range jobs {
		log.Printf("Worker %d received a message: %s", id, msg.Body)

		// Process the message
		var viewUpdate ViewUpdateMessage
		if err := json.Unmarshal(msg.Body, &viewUpdate); err != nil {
			log.Printf("Worker %d: Error unmarshaling message: %v", id, err)
			msg.Nack(false, true) // Nack and requeue
			continue
		}

		// Create paste view object
		pasteView := &model.PasteViews{
			PasteID:  viewUpdate.PasteID,
			ViewedAt: viewUpdate.ViewedAt,
		}

		// Update the database
		if err := repository.AddPasteViewDirect(pasteView); err != nil {
			log.Printf("Worker %d: Error adding paste view: %v", id, err)
			msg.Nack(false, true) // Nack and requeue
			continue
		}

		// Acknowledge the message
		msg.Ack(false)
		log.Printf("Worker %d: Successfully processed view update for paste %s", id, viewUpdate.PasteID)
	}

	log.Printf("Worker %d shutting down", id)
}
