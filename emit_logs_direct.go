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
	args := os.Args
	if len(args) == 2 || args[2] == "" {
		return "Default Message"
	}
	return strings.Join(args[2:], " ")
}

func severityCommandLineArgs() string {
	args := os.Args
	if len(args) == 1 || args[1] == "" {
		return "info"
	}
	return args[1]
}

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to create Channel")
	defer ch.Close()

	exchangeName := "logs_direct"
	err = ch.ExchangeDeclare(exchangeName, "direct", true, false, false, false, nil)
	failOnError(err, "Failed to declare Exchange")

	severity := severityCommandLineArgs()
	body := commandLineArgs()
	message := amqp.Publishing{
		ContentType: "text/plain",
		Body:        []byte(body),
	}
	err = ch.Publish(exchangeName, severity, false, false, message)
	failOnError(err, "Failed to publish message")

	log.Printf("Sent message: %s", body)
}
