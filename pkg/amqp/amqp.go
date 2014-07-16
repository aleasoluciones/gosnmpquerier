package amqp

import (
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

func Publish(amqpURI, exchange, routingKey, body string, reliable bool) error {

	connection, err := amqp.Dial(amqpURI)
	if err != nil {
		return fmt.Errorf("Dial: %s", err)
	}
	defer connection.Close()

	channel, err := connection.Channel()
	if err != nil {
		return fmt.Errorf("Channel: %s", err)
	}

	// Reliable publisher confirms require confirm.select support from the
	// connection.
	if reliable {
		log.Printf("enabling publishing confirms.")
		if err := channel.Confirm(false); err != nil {
			return fmt.Errorf("Channel could not be put into confirm mode: %s", err)
		}

		ack, nack := channel.NotifyConfirm(make(chan uint64, 1), make(chan uint64, 1))

		defer confirmOne(ack, nack)
	}

	if err = channel.Publish(
		exchange,   // publish to an exchange
		routingKey, // routing to 0 or more queues
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			Headers:         amqp.Table{},
			ContentType:     "text/json",
			ContentEncoding: "",
			Body:            []byte(body),
			DeliveryMode:    amqp.Transient, // 1=non-persistent, 2=persistent
			Priority:        0,              // 0-9
			// a bunch of application/implementation-specific fields
		},
	); err != nil {
		return fmt.Errorf("Exchange Publish: %s", err)
	}

	return nil
}

// One would typically keep a channel of publishings, a sequence number, and a
// set of unacknowledged sequence numbers and loop until the publishing channel
// is closed.
func confirmOne(ack, nack chan uint64) {
	log.Printf("waiting for confirmation of one publishing")

	select {
	case tag := <-ack:
		log.Printf("confirmed delivery with delivery tag: %d", tag)
	case tag := <-nack:
		log.Printf("failed delivery of delivery tag: %d", tag)
	}
}

func ReadJsonQueriesFromAmqp(amqpUri, exchangeName, queueName, bindingKey string) chan string {
	jsonQueries := make(chan string)

	go func() {
		conn, err := amqp.Dial(amqpUri)
		if err != nil {
			fmt.Println("Error dialing", err, amqpUri)
		}
		channel, err := conn.Channel()
		if err != nil {
			fmt.Println("Error geting channel", err)
		}
		err = channel.QueueBind(
			queueName,    // name of the queue
			bindingKey,   // bindingKey
			exchangeName, // sourceExchange
			false,        // noWait
			nil,          // arg
		)
		if err != nil {
			fmt.Println("Error queue binding", err, "queue", queueName, "bindingKey", bindingKey)
		}

		deliveries, err := channel.Consume(
			queueName, // name
			"tag-efa", // consumerTag,
			false,     // noAck
			false,     // exclusive
			false,     // noLocal
			false,     // noWait
			nil,       // arguments
		)
		if err != nil {
			fmt.Println("Error consuming", err, "queue", queueName)
		}
		for delivery := range deliveries {
			jsonQueries <- string(delivery.Body)
			delivery.Ack(false)

		}
	}()
	return jsonQueries
}
