package main

import (
	"log"

	"github.com/streadway/amqp"
)

type pub struct {
	RbbtChan    *amqp.Channel
	Queue       string
	Exchange    string
	PublishData amqp.Publishing
}

type PublishRbbt interface {
	Pub()
}

func NewPub(rbbtCh *amqp.Channel, queue, exchange string, pubData amqp.Publishing) PublishRbbt {
	return pub{RbbtChan: rbbtCh, Queue: queue, Exchange: exchange, PublishData: pubData}
}

func getRbbtChanConn(url string) (*amqp.Channel, error) {

	conn, err := amqp.Dial(url)

	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()

	if err != nil {
		return nil, err
	}

	return ch, nil
}

func (p pub) Pub() {

	p.RbbtChan.ExchangeDeclare(p.Exchange, "topic", true, false, false, false, nil)

	if err := p.RbbtChan.Publish(p.Exchange, "random-key", false, false, p.PublishData); err != nil {
		panic(err)
	}

	queue, err := p.RbbtChan.QueueDeclare(p.Queue, true, false, false, false, nil)

	log.Printf("\nQueue status: %v\n", queue)
	if err != nil {
		panic(err)
	}

	if err := p.RbbtChan.QueueBind(p.Queue, "#", p.Exchange, false, nil); err != nil {
		panic(err)
	}

	log.Println("Success to send data to queue :)")
}

func main() {

	url := "amqp://user:1234@localhost:5672"

	conn, err := amqp.Dial(url)

	if err != nil {
		panic("could not connect -> " + err.Error())
	}

	channel, err := conn.Channel()

	defer channel.Close()
	if err != nil {

		panic(err)
	}

	newPub := NewPub(channel, "TestQueue", "myExchange", amqp.Publishing{Body: []byte("Hello World from other app")})

	newPub.Pub()
	defer conn.Close()
}

//func Pub(chan []byte body,
