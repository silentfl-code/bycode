package main

import (
	"encoding/json"
	"fmt"
	"github.com/streadway/amqp"
	"os"
	"os/signal"
	"sync"
)

const (
	C           = 8 //count of consumers/producers/parsers
	RabbitMQUrl = "amqp://guest:guest@localhost:5672/"
	queue_limit = 200 //for every
)

type Message struct {
	Name string `json:"name"`
	Type string `json:"type"`
	Url  string `json:"url"`
}

type ParserFunc func(msg amqp.Delivery) bool

var parsers = map[string]ParserFunc{
	"sample": ParserSample,
}

var (
	conn         *amqp.Connection
	ch           *amqp.Channel
	queue_parser = make(chan amqp.Delivery, queue_limit)
	queue_writer = make(chan []byte, queue_limit)
	wg           sync.WaitGroup
	quit         = make(chan struct{}, 3)
)

func init() {
	conn, _ = amqp.Dial(RabbitMQUrl)
	ch, _ = conn.Channel()
	ch.QueueDeclare("test", false, false, false, true, nil)
}

func Writer() {
	for {
		select {
		case data := <-queue_writer:
			ch.Publish("", "test", false, false,
				amqp.Publishing{
					ContentType: "text/json",
					Body:        data,
				},
			)
		case <-quit:
			wg.Done()
		}
	}
}

func Reader() {
	messages, _ := ch.Consume("test", "", false, false, false, true, nil)
	for {
		select {
		case msg := <-messages:
			queue_parser <- msg
		case <-quit:
			wg.Done()
		}
	}
}

func Parser() {
	var v Message
	for {
		select {
		case msg := <-queue_parser:
			_ = json.Unmarshal(msg.Body, &v)
			parsers[v.Name](msg)
		case <-quit:
			wg.Done()
		}
	}
}

func main() {
	wg.Add(C * 3)

	for i := 0; i < C; i++ {
		go Reader()
		go Parser()
		go Writer()
	}

	waitCtrlC()

	for i := 0; i < C*3; i++ {
		quit <- struct{}{}
	}
	wg.Wait()
}

func waitCtrlC() {
	fmt.Println("Ctrl+C for exit")
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
}
