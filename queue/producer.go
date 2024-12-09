// queue/producer.go
package queue

import (
	"log"
	"github.com/streadway/amqp"
)

var conn *amqp.Connection
var channel *amqp.Channel

func InitRabbitMQ(url string) {
	var err error
	conn, err = amqp.Dial(url)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}

	channel, err = conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %v", err)
	}
}

func AddToQueue(images []string) {
	for _, img := range images {
		err := channel.Publish(
			"",           // exchange
			"image_queue", // routing key
			false,        // mandatory
			false,        // immediate
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        []byte(img),
			},
		)
		if err != nil {
			log.Printf("Failed to enqueue image: %v", err)
		}
	}
}
