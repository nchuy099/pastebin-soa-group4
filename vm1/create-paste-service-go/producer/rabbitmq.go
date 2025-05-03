package producer

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"time"

	"create-paste-service/model"

	"github.com/rabbitmq/amqp091-go"
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
	RabbitMQConn    *amqp091.Connection
	RabbitMQChannel *amqp091.Channel
)

// InitRabbitMQ initializes the RabbitMQ connection and channel
func InitRabbitMQ() error {
	rabbitURL := os.Getenv("RABBITMQ_URL")
	if rabbitURL == "" {
		rabbitURL = "amqp://guest:guest@localhost:5672/"
	}

	var err error
	RabbitMQConn, err = amqp091.Dial(rabbitURL)
	if err != nil {
		log.Printf("Failed to connect to RabbitMQ: %v", err)
		return err
	}

	RabbitMQChannel, err = RabbitMQConn.Channel()
	if err != nil {
		log.Printf("Failed to open a channel: %v", err)
		RabbitMQConn.Close()
		return err
	}

	// Declare a direct exchange
	err = RabbitMQChannel.ExchangeDeclare(
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

	log.Println("RabbitMQ connection established successfully")
	return nil
}

// CloseRabbitMQ closes the RabbitMQ connection and channel
func CloseRabbitMQ() {
	if RabbitMQChannel != nil {
		RabbitMQChannel.Close()
	}
	if RabbitMQConn != nil {
		RabbitMQConn.Close()
	}
}

// PublishPaste publishes a paste to the RabbitMQ queue
func PublishPaste(paste *model.Paste) error {
	if RabbitMQChannel == nil {
		return nil // Silently ignore if RabbitMQ is not connected
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Convert paste to message
	message := paste

	// Convert message to JSON
	body, err := json.Marshal(message)
	if err != nil {
		log.Printf("Error marshaling message: %v", err)
		return err
	}

	// Publish message to the exchange
	return RabbitMQChannel.PublishWithContext(
		ctx,
		ExchangeName, // exchange
		RoutingKey,   // routing key
		false,        // mandatory
		false,        // immediate
		amqp091.Publishing{
			ContentType:  "application/json",
			DeliveryMode: amqp091.Persistent, // Make message persistent
			Body:         body,
		},
	)
}
