package main

import (
	"fmt"

	amqp091 "github.com/rabbitmq/amqp091-go"
)

func main() {

	conn, err := amqp091.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		panic(err)
	}

	ch, err := conn.Channel()
	if err != nil {
		panic(err)
	}

	for i := 0; i < 100000; i++ {
		err := ch.Publish("", "queue", false, false, amqp091.Publishing{
			DeliveryMode: 0,
			ContentType:  "text/plain",
			Body:         []byte(fmt.Sprintf("%d", i)),
		})
		if err != nil {
			panic(err)
		}
	}

	if err := ch.Close(); err != nil {
		panic(err)
	}
	if err := conn.Close(); err != nil {
		panic(err)
	}
	fmt.Println("terminate")
}
