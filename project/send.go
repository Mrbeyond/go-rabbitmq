package main

import (
	"context"
	"log"
	"os"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func Sender() {
	// conn, err := amqp.Dial(os.Getenv("RABBITMQ_URL"))
	log.Println("\n\n\n\n ", os.Getenv("RABBITMQ_URL"), "\n\n\n\n ")
	conn, err := amqp.Dial(os.Getenv("RABBITMQ_URL"))
	FailOnError(err, "Cannot open rabbitMQ connection  from Sender \n\n ")
	defer conn.Close()

	ch, err := conn.Channel()
	FailOnError(err, "Cannot open rabbitMQ channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"hello", // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	FailOnError(err, "Cannot do QueueDeclare on rabbitMQ channel")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	log.Println(q.Name)
	body := "Helllo from sender"

	timer := time.NewTicker(2 * time.Second)

	for {
		select {
		case <-timer.C:
			err = ch.PublishWithContext(
				ctx,
				"",     // exchange
				q.Name, // routing key
				false,  // mandatory
				false,  // immediate
				amqp.Publishing{
					ContentType: "text/plain",
					Body:        []byte(body),
				},
			)

			FailOnError(err, "Failed to publish a message")
			log.Printf(" [x] Sent %s\n", body)

		case <-time.After(20 * time.Second):
			timer.Stop()
			return
		}
	}
}
