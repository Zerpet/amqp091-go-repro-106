package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	amqp091 "github.com/rabbitmq/amqp091-go"
)

var stop bool

func init() {
	flag.BoolVar(&stop, "stop", false, "Stop the program after 1 second")
	flag.Parse()
}

func main() {
	conn, err := amqp091.DialConfig("amqp://guest:guest@localhost:5672/", amqp091.Config{
		Heartbeat: time.Second * 30,
	})
	if err != nil {
		panic(err)
	}

	ch, err := conn.Channel()
	if err != nil {
		panic(err)
	}

	if err := ch.Qos(200, 0, false); err != nil {
		panic(err)
	}

	queueArgs := make(amqp091.Table)
	queueArgs["x-single-active-consumer"] = true
	// QueueDeclare(name string, durable, autoDelete, exclusive, noWait bool, args Table)
	queue, err := ch.QueueDeclare("queue",
		true,      // durable
		false,     // auto delete
		false,     //exclusive
		false,     //noWait
		queueArgs, // queue args
	)
	if err != nil {
		panic(err)
	}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	if stop {
		go func(c chan os.Signal) {
			<-time.After(time.Second*3)
			fmt.Println("Stopping by sending interrupt signal")
			c <- os.Interrupt
		}(sigs)
	}
	msgs, err := ch.Consume(queue.Name, "consumer", false, false, false, false, nil)
	if err != nil {
		panic(err)
	}
	// d := make(chan bool)
	go func() {
		for msg := range msgs {
			fmt.Println(string(msg.Body))
			_ = msg.Ack(true)
		}
		// d <- true
		// for {
		// 	select {
		// 		case msg := <-msgs:
		// 			fmt.Println(string(msg.Body))
		// 			_ = msg.Ack(true)
		// 		case <-d:
    //
		// 	}
		// }
	}()

	<-sigs

	if err := ch.Cancel("consumer", false); err != nil {
		panic(err)
	}
	fmt.Println("cancel consume")
	if err := ch.Close(); err != nil {
		panic(err)
	}
	if err := conn.Close(); err != nil {
		panic(err)
	}
	fmt.Println("terminate")
}
