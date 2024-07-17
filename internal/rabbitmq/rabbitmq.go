package rabbitmq

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"log"
	"rabbitmq/internal/models"
	"rabbitmq/internal/storage"

	"github.com/streadway/amqp"
)

func SendToQueue(body []byte, contentType string) {
	ch, err := storage.Conn()
	if err != nil {
		log.Fatal(err)
	}

	q, err := ch.QueueDeclare(
		"task_queue",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to declare a queue: %s", err)
	}

	msg := amqp.Publishing{
		ContentType: contentType,
		Body:        body,
	}

	if err = ch.Publish(
		"",
		q.Name,
		false,
		false,
		msg,
	); err != nil {
		log.Fatalf("Failed to publish a message: %s", err)
	}

}

func StartConsumer() {
	ch, err := storage.Conn()
	if err != nil {
		log.Fatal(err)
	}
	q, err := ch.QueueDeclare(
		"task_queue",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to declare a queue: %s", err)
	}
	if err != nil {
		log.Fatalf("Failed to declare queue: %s", err)
	}

	msgs, err := ch.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to register a consumer: %s", err)
	}
	forever := make(chan bool)
	go func() {
		for d := range msgs {
			if d.ContentType == "application/json" {
				var user models.User
				if err := json.Unmarshal(d.Body, &user); err != nil {
					log.Printf("Error decoding JSON: %s", err)
				} else {
					fmt.Println("Received JSON:", user)
				}
			} else if d.ContentType == "application/xml" {
				var xmlUser models.User
				if err := xml.Unmarshal(d.Body, &xmlUser); err != nil {
					log.Printf("Error decoding XML: %s", err)
				} else {
					fmt.Println("Received XML:", xmlUser)
				}
			}
		}
	}()

	log.Printf("Waiting  for messages. To exit press CTRL+C")
	<-forever
}
