package main

import (
	"fmt"
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
	connectionUrl := fmt.Sprintf("amqp://%s:%s@%s:%s/", os.Getenv("RMQ_USER"), os.Getenv("RMQ_PASS"), os.Getenv("RMQ_HOST"), os.Getenv("RMQ_PORT"))
	// establish a RMQ connection
	conn, err := amqp.Dial(connectionUrl)
	failOnError(err, "Failed to Connect")
	defer conn.Close()

	// channel to the connection made to start interaction with RMQ
	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	// Make a queue via the interaction
	// If we send message to a non-existing queue, RMQ will just drop the message
	// QueueDeclare is idempotent - only 1 is created even if we declare it multiple times
	// Marking the queue as durable so that it is not lost even if RMQ restarts
	q, err := ch.QueueDeclare("TestQueue", true, false, false, false, nil)
	failOnError(err, "Failed to declare a queue")

	body := commandLineArgs()
	// Publish a message over to the queue
	// Marking messages as persistent doesn't fully guarantee that a message won't be lost. There exists a short
	// time when message is received by RMQ and hasn't been saved to the disk. Also it may not even store as a
	// persistent storage, but use cache as its mechanism.
	msg := amqp.Publishing{
		ContentType:  "text/plain",
		Body:         []byte(body),
		DeliveryMode: amqp.Persistent,
	}
	/*
		messages can never be sent directly to a queue, it always needs to go through an exchange.
		Default Exchange => Empty String
	*/
	err = ch.Publish("", q.Name, false, false, msg)
	failOnError(err, "Failed to publish a message")
	log.Printf(" [x] Sent %s", body)
}
