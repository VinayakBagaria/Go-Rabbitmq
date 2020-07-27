package routing

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

func severityCommandLineArgs() []string {
	args := os.Args
	if len(args) == 1 {
		log.Fatal("Expected either warn info error as 2nd argument on")
		os.Exit(0)
	}
	return args[1:]
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

	q, err := ch.QueueDeclare("", false, false, true, false, nil)
	failOnError(err, "Failed to declare Queue")

	listenerLogTypes := severityCommandLineArgs()
	for _, s := range listenerLogTypes {
		log.Printf("Binding Queue %s to exchange %s on direct type", q.Name, s)
		err = ch.QueueBind(q.Name, s, exchangeName, false, nil)
		failOnError(err, "Failed to bind a Queue")
	}

	msgs, err := ch.Consume(q.Name, "", true, false, false, false, nil)
	failOnError(err, "Failed to register a Consumer")

	forever := make(chan bool)
	go func(){
		for d := range msgs {
			log.Printf("[x] %s", d.Body)
		}
	}()
	log.Printf(" [*] Waiting for logs. To exit press CTRL+C")
	<-forever
}
