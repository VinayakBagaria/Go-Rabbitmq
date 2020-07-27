package main

import (
	"fmt"
	"log"
	"github.com/streadway/amqp"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatal("%s: %s", msg, err)
		panic(err)
	}
}

func main() {
	fmt.Println("Consumer App")

	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to Connect")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to Open a Channel")
	defer ch.Close()

	q, err := ch.QueueDeclare("TestQueue", false, false,false, false, nil)
	failOnError(err, "Failed to Declare a Queue")

	// assuming here queue is made before and hence no ch.QueueDeclare()
	msgs, err := ch.Consume(q.Name, "", true, false, false, false, nil)

	forever := make(chan bool)
	go func(){
		for d := range msgs {
			fmt.Printf("Received Message: %s\n", d.Body)
		}
	}()

	fmt.Println("Successfully Connected to RMQ")
	fmt.Println(" [*] - waiting for msgs")
	// since no one is sending to channel, this is forever blocking to get a value
	<-forever
}
