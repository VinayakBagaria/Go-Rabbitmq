package main

import (
	"bytes"
	"github.com/streadway/amqp"
	"log"
	"time"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatal("%s: %s", msg, err)
		panic(err)
	}
}

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to Connect")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to Open a Channel")
	defer ch.Close()

	q, err := ch.QueueDeclare("TestQueue", true, false, false, false, nil)
	failOnError(err, "Failed to Declare a Queue")

	// Do not auto-ack once the task is consumed, we will manually do it so that RMQ doesn't delete the task
	// once consumed
	msgs, err := ch.Consume(q.Name, "", false, false, false, false, nil)

	forever := make(chan bool)
	go func() {
		for d := range msgs {
			// message here will byte slice
			message := d.Body
			log.Printf("Received a message: %s", message)
			// count of dots in message bytes via a dot byte slice
			dot_count := bytes.Count(message, []byte("."))
			t := time.Duration(dot_count)
			time.Sleep(t * time.Second)
			log.Printf("[x] Done with %s after %v seconds", message, dot_count)
			// acknowledge that the task is completed
			d.Ack(false)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	// since no one is sending to channel, this is forever blocking to get a value
	<-forever
}
