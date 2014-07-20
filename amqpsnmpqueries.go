package main

import (
	"flag"
	"fmt"

	"github.com/eferro/go-snmpqueries/pkg/amqp"
	"github.com/eferro/go-snmpqueries/pkg/snmpquery"
)

const (
	CONTENTION = 4
)

func publishResults(processed chan snmpquery.Query, amqpUri, exchangeName, routingKey string) {
	for query := range processed {
		jsonText, err := snmpquery.ToJson(&query)
		if err != nil {
			fmt.Println("Error converting to json:", err, query)
			continue
		}
		amqp.Publish(amqpUri, exchangeName, routingKey, jsonText, false)
		fmt.Println("Publicado:", jsonText)
	}
}

func dispatchFromAmqpToInput(input chan snmpquery.Query, amqpUri, sourceExchange, queueName, bindingKey string) {
	for jsonQuery := range amqp.ReadJsonQueriesFromAmqp(amqpUri, sourceExchange, queueName, bindingKey) {
		query, err := snmpquery.FromJson(jsonQuery)
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

	querier := snmpquery.New(CONTENTION)
	go dispatchFromAmqpToInput(querier.Input, *amqpUri, *sourceExchange, *queue, *bindingKey)
	publishResults(querier.Output, *amqpUri, *destinationExchange, *publishKey)
}
