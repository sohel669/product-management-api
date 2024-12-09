package services

import (
	"bytes"
	"database/sql"
	"image"
	"image/jpeg"
	"log"
	"net/http"
	"os"

	"github.com/streadway/amqp"
	_ "github.com/lib/pq"
)

var db *sql.DB

func InitializeImageProcessing(dsn string) {
	var err error
	db, err = sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
}

// ConsumeImageQueue processes images from RabbitMQ
func ConsumeImageQueue(rabbitMQURL, queueName string) {
	conn, err := amqp.Dial(rabbitMQURL)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %v", err)
	}
	defer ch.Close()

	msgs, err := ch.Consume(
		queueName, // queue
		"",        // consumer
		true,      // auto-ack
		false,     // exclusive
		false,     // no-local
		false,     // no-wait
		nil,       // args
	)
	if err != nil {
		log.Fatalf("Failed to register a consumer: %v", err)
	}

	for d := range msgs {
		processImage(string(d.Body))
	}
}

func processImage(imageURL string) {
	resp, err := http.Get(imageURL)
	if err != nil {
		log.Printf("Failed to download image: %v", err)
		return
	}
	defer resp.Body.Close()

	img, _, err := image.Decode(resp.Body)
	if err != nil {
		log.Printf("Failed to decode image: %v", err)
		return
	}

	// Compress the image
	var compressed bytes.Buffer
	err = jpeg.Encode(&compressed, img, &jpeg.Options{Quality: 50})
	if err != nil {
		log.Printf("Failed to compress image: %v", err)
		return
	}

	// Save the compressed image
	fileName := "/tmp/compressed_image.jpg"
	err = os.WriteFile(fileName, compressed.Bytes(), 0644)
	if err != nil {
		log.Printf("Failed to save compressed image: %v", err)
		return
	}

	// Update database with compressed image
	query := `UPDATE products SET compressed_product_images = array_append(compressed_product_images, $1) WHERE $2 = ANY(product_images)`
	_, err = db.Exec(query, fileName, imageURL)
	if err != nil {
		log.Printf("Failed to update database: %v", err)
		return
	}

	log.Printf("Image processed and saved: %s", fileName)
}
