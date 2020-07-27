package main

import (
	"fmt"
	"github.com/streadway/amqp"
)

func main() {
	fmt.Println("Go RMQ")

	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	defer conn.Close()

	fmt.Println("Successfully connected to RMQ")

	ch, err := conn.Channel()
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	defer ch.Close()

	q, err := ch.QueueDeclare("TestQueue", false, false, false, false, nil)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	fmt.Println(q)

	msg := amqp.Publishing{
		ContentType: "text/plain",
		Body:        []byte("Hello World"),
	}
	err = ch.Publish("", "TestQueue", false, false, msg)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	fmt.Println("Successfully published msg to queue")
}
