package main

import (
	"github.com/streadway/amqp"
	"log"
	"os"
	"strings"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatal("%s: %s", msg, err)
		panic(err)
	}
}

func commandLineArgs() string {
	//var s string
	/*
	Command executed like "go run main.go some text here"
	os.Args will return an array of length 4 -> [exec-path-to-main-go some text here]
	We want all string from and after index 1
	 */
	args := os.Args
	// no argument passed or empty passed
	if len(args) < 2 || args[1] == "" {
		return "Default Message"
	}
	return strings.Join(args[1:], " ")
}

func main() {
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

	body := commandLineArgs()
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
