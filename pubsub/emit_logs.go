package pubsub

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
	if len(args) < 2 || args[1] == "" {
		return "Default Message"
	}
	return strings.Join(args[1:], " ")
}

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to Connect")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	// declaring an exchange "logs" of type fanout(pub/sub)
	exchangeName := "logs"
	err = ch.ExchangeDeclare(exchangeName, "fanout", true, false, false, false, nil)
	failOnError(err, "Failed to create an exchange")

	body := commandLineArgs()
	msg := amqp.Publishing{
		ContentType: "text/plain",
		Body:        []byte(body),
	}
	// message is lost if no consumer has a queue is bound to this exchange
	err = ch.Publish(exchangeName, "", false, false, msg)
	failOnError(err, "Failed to publish a message")
	log.Printf(" [x] Sent %s", body)
}
