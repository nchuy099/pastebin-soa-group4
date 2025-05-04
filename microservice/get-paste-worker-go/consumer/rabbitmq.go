package consumer

import (
	"log"
	"os"
	"time"

	"github.com/rabbitmq/amqp091-go"
)

const (
	ExchangeName = "view_exchange"
	QueueName    = "view_queue"
	RoutingKey   = "view.add"
)

var (
	// Connection is the RabbitMQ connection
	Connection *amqp091.Connection
	// Channel is the RabbitMQ channel
	Channel *amqp091.Channel
)

// ViewUpdateMessage structure for RabbitMQ messages
type ViewUpdateMessage struct {
	PasteID    string    `json:"paste_id"`
	ViewedAt   time.Time `json:"viewed_at"`
	RemoteAddr string    `json:"remote_addr,omitempty"`
	UserAgent  string    `json:"user_agent,omitempty"`
}

// InitRabbitMQ initializes the RabbitMQ connection and channel
func InitRabbitMQ() error {
	// Get RabbitMQ URL from environment
	rabbitURL := os.Getenv("RABBITMQ_URL")
	if rabbitURL == "" {
		rabbitURL = "amqp://guest:guest@localhost:5672/"
	}

	var err error

	// Connect to RabbitMQ
	Connection, err = amqp091.Dial(rabbitURL)
	if err != nil {
		log.Printf("Failed to connect to RabbitMQ: %v", err)
		return err
	}

	// Create a channel
	Channel, err = Connection.Channel()
	if err != nil {
		log.Printf("Failed to open a channel: %v", err)
		return err
	}

	// Declare exchange
	err = Channel.ExchangeDeclare(
		ExchangeName, // name
		"direct",     // type
		true,         // durable
		false,        // auto-deleted
		false,        // internal
		false,        // no-wait
		nil,          // arguments
	)
	if err != nil {
		log.Printf("Failed to declare an exchange: %v", err)
		return err
	}

	// Declare queue
	_, err = Channel.QueueDeclare(
		QueueName, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		log.Printf("Failed to declare a queue: %v", err)
		return err
	}

	// Bind queue to exchange
	err = Channel.QueueBind(
		QueueName,    // queue name
		RoutingKey,   // routing key
		ExchangeName, // exchange
		false,
		nil,
	)
	if err != nil {
		log.Printf("Failed to bind a queue: %v", err)
		return err
	}

	log.Println("Successfully connected to RabbitMQ")
	return nil
}

// CloseRabbitMQ closes the RabbitMQ connection and channel
func CloseRabbitMQ() {
	if Channel != nil {
		Channel.Close()
	}
	if Connection != nil {
		Connection.Close()
	}
}

// ConsumeMessages consumes messages from the queue
func ConsumeMessages() (<-chan amqp091.Delivery, error) {
	// Start consuming messages
	msgs, err := Channel.Consume(
		QueueName, // queue
		"",        // consumer
		false,     // auto-ack
		false,     // exclusive
		false,     // no-local
		false,     // no-wait
		nil,       // args
	)
	if err != nil {
		log.Printf("Failed to register a consumer: %v", err)
		return nil, err
	}

	log.Println("Successfully registered as a consumer")
	return msgs, nil
}
