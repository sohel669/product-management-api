package queue

import (
	"log"

	"project/database"
	"project/services"
	"github.com/streadway/amqp"
)

func StartConsumer(queueName string) {
	msgs, err := channel.Consume(
		queueName, // queue
		"",        // consumer tag
		true,      // auto-ack
		false,     // exclusive
		false,     // no-local
		false,     // no-wait
		nil,       // args
	)
	if err != nil {
		log.Fatalf("Failed to register consumer: %v", err)
	}

	go func() {
		for msg := range msgs {
			imageURL := string(msg.Body)
			compressedImage := services.ProcessImage(imageURL)

			query := `UPDATE products SET compressed_product_images = array_append(compressed_product_images, $1) WHERE $2 = ANY(product_images)`
			_, err := database.DB.Exec(query, compressedImage, imageURL)
			if err != nil {
				log.Printf("Failed to update product with compressed image: %v", err)
			}
		}
	}()
}
