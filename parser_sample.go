package main

import (
	"encoding/json"
	"fmt"
	"github.com/moovweb/gokogiri"
	"github.com/streadway/amqp"
	"io/ioutil"
	"net/http"
)

var ParserSample ParserFunc = func(msg amqp.Delivery) bool {
	// затычка, поскольку нигде нет обработки ошибок
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("parser sample: recovered")
		}
	}()
	var v Message
	json.Unmarshal(msg.Body, &v)
	r, _ := http.Get(v.Url)
	defer r.Body.Close()
	body, _ := ioutil.ReadAll(r.Body)
	doc, _ := gokogiri.ParseHtml(body)

	switch v.Type {
	case "start":
		_, _ = doc.Search("//title")
		//title, _ := doc.Search("//title")	//uncomment this if print need
		//fmt.Println(title[0].Content())

		urls, _ := doc.Search("//a/@href")
		for _, url := range urls {
			data, _ := json.Marshal(Message{
				v.Name,
				"level1",
				url.String(),
			})
			queue_writer <- data
		}
		msg.Ack(false)
	case "level1":
		_, _ = doc.Search("//title")
		//title, _ := doc.Search("//title")	//uncomment this if print need
		//fmt.Println("level1: ", title)
		msg.Ack(false)
	default:
		fmt.Println("parser sample: not supported type")
	}
	return true
}
