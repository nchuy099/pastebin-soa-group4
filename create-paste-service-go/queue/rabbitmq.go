package queue

import (
	"encoding/json"
	"log"
	"os"
	"time"

	"create-paste-service-go/model"

	"github.com/streadway/amqp"
)

const (
	// ExchangeName is the name of the direct exchange
	ExchangeName = "paste_exchange"
	// QueueName is the name of the queue
	QueueName = "paste_queue"
	// RoutingKey is the routing key for paste messages
	RoutingKey = "paste.create"
)

var (
	// Connection is the RabbitMQ connection
	Connection *amqp.Connection
	// Channel is the RabbitMQ channel
	Channel *amqp.Channel
)

// PasteMessage represents a message to be sent to the queue
type PasteMessage struct {
	ID         string     `json:"id"`
	Content    string     `json:"content"`
	Title      string     `json:"title"`
	Language   string     `json:"language"`
	CreatedAt  time.Time  `json:"created_at"`
	ExpiresAt  *time.Time `json:"expires_at,omitempty"`
	Views      int        `json:"views"`
	Visibility string     `json:"visibility"`
}

// InitRabbitMQ initializes the RabbitMQ connection and channel
func InitRabbitMQ() error {
	// Get RabbitMQ URL from environment
	rabbitURL := os.Getenv("RABBITMQ_URL")
	if rabbitURL == "" {
		rabbitURL = "amqp://guest:guest@rabbitmq:5672/"
	}

	var err error

	// Connect to RabbitMQ
	Connection, err = amqp.Dial(rabbitURL)
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

	// Set QoS for better load distribution
	err = Channel.Qos(
		1,     // prefetch count - only process one message at a time
		0,     // prefetch size - no specific size limit
		false, // global - apply to this channel only
	)
	if err != nil {
		log.Printf("Failed to set QoS: %v", err)
		return err
	}

	// Declare a direct exchange
	err = Channel.ExchangeDeclare(
		ExchangeName, // name
		"direct",     // type - using direct exchange as specified
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

	// Declare a queue
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

	// Bind the queue to the exchange with the routing key
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

// PublishPaste publishes a paste to the RabbitMQ queue
func PublishPaste(paste *model.Paste) error {
	// Convert paste to message
	message := PasteMessage{
		ID:         paste.ID,
		Content:    paste.Content,
		Title:      paste.Title,
		Language:   paste.Language,
		CreatedAt:  paste.CreatedAt,
		ExpiresAt:  paste.ExpiresAt,
		Visibility: string(paste.Visibility),
	}

	// Convert message to JSON
	body, err := json.Marshal(message)
	if err != nil {
		log.Printf("Error marshaling message: %v", err)
		return err
	}

	// Publish message to the exchange
	err = Channel.Publish(
		ExchangeName, // exchange
		RoutingKey,   // routing key
		false,        // mandatory
		false,        // immediate
		amqp.Publishing{
			ContentType:  "application/json",
			DeliveryMode: amqp.Persistent, // Make message persistent
			Body:         body,
		})
	if err != nil {
		log.Printf("Error publishing message: %v", err)
		return err
	}

	log.Printf("Successfully published message for paste ID: %s", paste.ID)
	return nil
}

// ConsumeMessages consumes messages from the queue
func ConsumeMessages() (<-chan amqp.Delivery, error) {
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
