package main

import (
	"log"

	"github.com/streadway/amqp"
)

func main() {

	conn, err := amqp.Dial("amqp://user:1234@localhost:5672")

	if err != nil {
		panic(err)
	}

	ch, err := conn.Channel()

	if err != nil {
		panic(err)
	}

	defer ch.Close()

	msgs, err := ch.Consume("TestQueue", "myExchange", true, false, false, false, nil)

	forever := make(chan bool)

	go func() {

		for m := range msgs {
			log.Printf("Receive msg: %v", string(m.Body))
		}
	}()

	<-forever
}
