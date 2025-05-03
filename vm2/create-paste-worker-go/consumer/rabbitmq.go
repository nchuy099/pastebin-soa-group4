package consumer

import (
	"log"
	"os"

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

	// Set QoS for better load distribution
	err = RabbitMQChannel.Qos(
		1,     // prefetch count - only process one message at a time
		0,     // prefetch size - no specific size limit
		false, // global - apply to this channel only
	)
	if err != nil {
		log.Printf("Failed to set QoS: %v", err)
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

	// Declare a queue
	_, err = RabbitMQChannel.QueueDeclare(
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
	err = RabbitMQChannel.QueueBind(
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

// ConsumeMessages consumes messages from the queue
func ConsumeMessages() (<-chan amqp091.Delivery, error) {
	// Start consuming messages
	msgs, err := RabbitMQChannel.Consume(
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
