package main

import (
	"fmt"
	"github.com/streadway/amqp"
	"log"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatal("%s: %s", msg, err)
		panic(err)
	}
}

func main() {
	fmt.Println("Go RMQ")

	// establish a RMQ connection
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to Connect")
	defer conn.Close()

	// channel to the connection made to start interaction with RMQ
	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	// make a queue via the interaction
	// if we send message to a non-existing queue, RMQ will just drop the message
	// QueueDeclare is idempotent - 1 is created even if we declare it multiple times
	q, err := ch.QueueDeclare("TestQueue", false, false, false, false, nil)
	failOnError(err, "Failed to declare a queue")
	fmt.Println(q)

	body := "Hello World"
	// publish a message over to the queue
	msg := amqp.Publishing{
		ContentType: "text/plain",
		Body:        []byte(body),
	}
	/*
	messages can never be sent directly to a queue, it always needs to go through an exchange.
	Default Exchange => Empty String
	 */
	err = ch.Publish("", q.Name, false, false, msg)
	failOnError(err, "Failed to publish a message")
	log.Printf(" [x] Sent %s", body)
}
