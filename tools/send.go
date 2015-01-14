package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"encoding/json"
	"flag"
	"github.com/streadway/amqp"
)

var (
	addr = flag.String("addr", "http://localhost:3000/%d", "url which will pushed to RabbitMQ")
	N    = flag.Int("N", 1, "count of messages")
)

func init() {
	flag.Parse()
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
		panic(fmt.Sprintf("%s: %s", msg, err))
	}
}

type Message struct {
	Name string `json:"name"`
	Type string `json:"type"`
	Url  string `json:"url"`
}

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"test", // name
		false,  // durable
		false,  // delete when usused
		false,  // exclusive
		true,   // no-wait
		nil,    // arguments
	)
	failOnError(err, "Failed to declare a queue")

	for i := 1; i <= *N; i++ {
		body, _ := json.Marshal(Message{"sample", "start", fmt.Sprintf(*addr, i)})
		err = ch.Publish(
			"",     // exchange
			q.Name, // routing key
			false,  // mandatory
			false,  // immediate
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        []byte(body),
			})
		failOnError(err, "Failed to publish a message")
	}
}
