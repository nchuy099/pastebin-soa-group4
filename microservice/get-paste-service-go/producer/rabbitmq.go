package producer

import (
	"context"
	"encoding/json"
	"get-paste-service/model"
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

	// Declare exchange
	err = RabbitMQChannel.ExchangeDeclare(
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

// PublishPasteView publishes a paste view message to RabbitMQ
func PublishPasteView(pasteView *model.PasteView) error {
	if RabbitMQChannel == nil {
		return nil // Silently ignore if RabbitMQ is not connected
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	message := pasteView

	body, err := json.Marshal(message)
	if err != nil {
		log.Printf("Error marshaling message: %v", err)
		return err
	}

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
