package pubsub

import (
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
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to Connect")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	exchange_name := "logs"
	err = ch.ExchangeDeclare(exchange_name, "fanout", true, false, false, false, nil)
	failOnError(err, "Failed to create an exchange")

	// Whenever we connect on this pub/sub mechanism, make a fresh/empty queue and let RMQ decide its name
	// via empty string
	// With exclusive set to true, this queue will be deleted once this consumer's connection is closed
	q, err := ch.QueueDeclare("", false, false, true, false, nil)
	failOnError(err, "Failed to declare a queue")

	// Relationship b/w Exchange and Queue is binding
	// We tell the logs exchange to send message to our newly created queue for this consumer
	err = ch.QueueBind(q.Name, "", exchange_name, false, nil)
	failOnError(err, "Failed to bind a queue")

	msgs, err := ch.Consume(q.Name, "", true, false, false, false, nil)
	failOnError(err, "Failed to register a consumer")

	forever := make(chan bool)
	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
		}
	}()

	log.Printf(" [*] Waiting for logs. To exit press CTRL+C")
	<-forever
}
