Simple multithreaded parser based RabbitMQ + Golang
===================================================

depedency:
---------
* rabbitmq
* go/gccgo
* some packages:
  - "github.com/streadway/amqp"
  - "github.com/moovweb/gokogiri"

compile tools:
--------------
	go build -o tools/send tools/send.go
	go build -o tools/http_server tools/http_server.go

compile this:
-------------
	go build *.go

before running ./main need:
* publish some messages (i.e. _./tools/send -N 100000_)
* started _./tools/http_server_
