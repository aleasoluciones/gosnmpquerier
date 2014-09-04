package main

import (
    "flag"
    "fmt"
    "log"

    "github.com/soniah/gosnmp"
    "github.com/streadway/amqp"
)

const (
    CONTENTION = 4
)

func publish(amqpURI, exchange, routingKey, body string, reliable bool) error {

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

func readJsonQueriesFromAmqp(amqpUri, exchangeName, queueName, bindingKey string) chan string {
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

func publishResults(processed chan gosnmpquerier.Query, amqpUri, exchangeName, routingKey string) {
    for query := range processed {
        jsonText, err := gosnmpquerier.ToJson(&query)
        if err != nil {
            fmt.Println("Error converting to json:", err, query)
            continue
        }
        publish(amqpUri, exchangeName, routingKey, jsonText, false)
        fmt.Println("Publicado:", jsonText)
    }
}

func dispatchFromAmqpToInput(input chan gosnmpquerier.Query, amqpUri, sourceExchange, queueName, bindingKey string) {
    for jsonQuery := range readJsonQueriesFromAmqp(amqpUri, sourceExchange, queueName, bindingKey) {
        query, err := gosnmpquerier.FromJson(jsonQuery)
        if err != nil {
            fmt.Println("ERROR", err, jsonQuery)
        } else {
            input <- *query
        }
    }
}

func main() {
    amqpUri := flag.String("amqp_uri", "amqp://guest:guest@localhost:5672/", "AMQP uri")

    queue := flag.String("queue", "queue1", "Queue name")
    sourceExchange := flag.String("src_exchange", "src_exchange1", "Durable AMQP exchange name")
    bindingKey := flag.String("binding_key", "", "AMQP routing key")

    destinationExchange := flag.String("dst_exchange", "dest_exchange1", "Durable AMQP exchange name")
    publishKey := flag.String("publish_key", "", "AMQP routing key")
    flag.Parse()

    querier := gosnmpquerier.NewAsyncQuerier(CONTENTION)
    go dispatchFromAmqpToInput(querier.Input, *amqpUri, *sourceExchange, *queue, *bindingKey)
    publishResults(querier.Output, *amqpUri, *destinationExchange, *publishKey)
}
