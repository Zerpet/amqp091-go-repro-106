# amqp091-go-repro-106
Repro code for rabbitmq/amqp091-go/issues/106

## How to run this code

You need 3 terminals and RabbitMQ running, listening in the default AMQP port 5672.

Prepare the following commands in different terminals:

```bash
go run consumer/consumer.go --stop
```

The above will exit after 3 seconds. Don't hit enter until you have all them ready :-)
Or change the code in this line: https://github.com/Zerpet/amqp091-go-repro-106/blob/f65a7987a759f25ee7602b3c4c3becd25a8f6320/consumer/consumer.go#L56

```bash
go run consumer/consumer.go
```

The above won't exit. It should be "ready" consumer. First one should be active.

```bash
go run producer/producer.go
```

The above will publish 100,000 messages and exit. Producer relies on Consumers to create the test queue.
