package queue

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"time"

	"github.com/rabbitmq/amqp091-go"
)

var (
	RabbitMQConn    *amqp091.Connection
	RabbitMQChannel *amqp091.Channel
)

const (
	ViewUpdateExchange = "view_updates"
	ViewUpdateQueue    = "view_updates_queue"
	ViewUpdateKey      = "view.update"
)

// ViewUpdateMessage represents the message structure for view updates
type ViewUpdateMessage struct {
	PasteID    string    `json:"paste_id"`
	ViewedAt   time.Time `json:"viewed_at"`
	RemoteAddr string    `json:"remote_addr,omitempty"`
	UserAgent  string    `json:"user_agent,omitempty"`
}

// InitRabbitMQ initializes the RabbitMQ connection and channel
func InitRabbitMQ() {
	rabbitURL := os.Getenv("RABBITMQ_URL")
	if rabbitURL == "" {
		rabbitURL = "amqp://guest:guest@localhost:5672/"
	}

	var err error
	RabbitMQConn, err = amqp091.Dial(rabbitURL)
	if err != nil {
		log.Printf("Failed to connect to RabbitMQ: %v", err)
		return
	}

	RabbitMQChannel, err = RabbitMQConn.Channel()
	if err != nil {
		log.Printf("Failed to open a channel: %v", err)
		RabbitMQConn.Close()
		return
	}

	// Declare exchange
	err = RabbitMQChannel.ExchangeDeclare(
		ViewUpdateExchange, // name
		"direct",           // type
		true,               // durable
		false,              // auto-deleted
		false,              // internal
		false,              // no-wait
		nil,                // arguments
	)
	if err != nil {
		log.Printf("Failed to declare an exchange: %v", err)
		return
	}

	log.Println("RabbitMQ connection established successfully")
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

// PublishViewUpdate publishes a view update message to RabbitMQ
func PublishViewUpdate(pasteID string, viewedAt time.Time, remoteAddr, userAgent string) error {
	if RabbitMQChannel == nil {
		return nil // Silently ignore if RabbitMQ is not connected
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	message := ViewUpdateMessage{
		PasteID:    pasteID,
		ViewedAt:   viewedAt,
		RemoteAddr: remoteAddr,
		UserAgent:  userAgent,
	}

	body, err := json.Marshal(message)
	if err != nil {
		return err
	}

	return RabbitMQChannel.PublishWithContext(
		ctx,
		ViewUpdateExchange, // exchange
		ViewUpdateKey,      // routing key
		false,              // mandatory
		false,              // immediate
		amqp091.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
}
