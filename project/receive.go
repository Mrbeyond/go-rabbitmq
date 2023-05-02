package main

import (
	"log"
	"os"

	amqp "github.com/rabbitmq/amqp091-go"
)

func Receiver() {

	log.Println("\n\n\n\n  From Receiver", os.Getenv("RABBITMQ_URL"), "\n\n\n\n ")
	conn, err := amqp.Dial(os.Getenv("RABBITMQ_URL"))
	FailOnError(err, "Cannot open rabbitMQ connection from Receiver \n\n ")
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

	msgs, err := ch.Consume(
		q.Name, //queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	FailOnError(err, "Cannot on rabbitMQ channel Consume")
	var forever chan struct{}

	go func() {
		for d := range msgs {
			log.Printf("\n\n\n Received a message: %s \n\n\n ", d.Body)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
