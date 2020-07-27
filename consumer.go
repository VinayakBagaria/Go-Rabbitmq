package main

import (
	"fmt"
	"github.com/streadway/amqp"
)

func main() {
	fmt.Println("Consumer App")

	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	defer ch.Close()

	msgs, err := ch.Consume("TestQueue", "", true, false, false, false, nil)

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
