package snmpquery

import (
	"fmt"
	"time"

	"github.com/soniah/gosnmp"
)

type OpSnmp int32

const (
	GET  = 0
	WALK = 1
)

type Query struct {
	Id          int
	Cmd         OpSnmp
	Community   string
	Oid         string
	Timeout     time.Duration
	Retries     int
	Destination string
	Response    []gosnmp.SnmpPDU
	Error       error
}

func Process(input chan Query, processed chan Query, conntention int) {
	fmt.Println("EFA DELETE")

	m := make(map[string]chan Query)

	for query := range input {
		_, exists := m[query.Destination]
		if exists == false {
			channel_tmp := make(chan Query, 10)
			m[query.Destination] = channel_tmp
			for i := 0; i < conntention; i++ {
				go processQueriesFromChannel(channel_tmp, processed)
			}
		}
		m[query.Destination] <- query
	}
}

func handleQuery(query *Query) {

	switch query.Cmd {
	case WALK:
		query.Response, query.Error = walk(query.Destination, query.Community, query.Oid, query.Timeout, query.Retries)
	case GET:
		query.Response, query.Error = get(query.Destination, query.Community, query.Oid, query.Timeout, query.Retries)
	}
}

func walk(destination, community, oid string, timeout time.Duration, retries int) ([]gosnmp.SnmpPDU, error) {
	conn := snmpConnection(destination, community, timeout, retries)
	err := conn.Connect()
	if err != nil {
		return nil, err
	}
	defer conn.Conn.Close()
	output := make(chan gosnmp.SnmpPDU)
	errChannel := make(chan error, 1)
	go doWalk(conn, oid, output, errChannel)

	result := []gosnmp.SnmpPDU{}
	for pdu := range output {
		result = append(result, pdu)
	}
	if len(errChannel) != 0 {
		err := <-errChannel
		return nil, err
	}
	return result, nil
}

func doWalk(conn gosnmp.GoSNMP, oid string, output chan gosnmp.SnmpPDU, errChannel chan error) {
	processPDU := func(pdu gosnmp.SnmpPDU) error {
		output <- pdu
		return nil
	}
	err := conn.BulkWalk(oid, processPDU)
	if err != nil {
		errChannel <- err
	}
	close(output)
}

func snmpConnection(destination, community string, timeout time.Duration, retries int) gosnmp.GoSNMP {
	return gosnmp.GoSNMP{
		Target:    destination,
		Port:      161,
		Community: community,
		Version:   gosnmp.Version2c,
		Timeout:   timeout,
		Retries:   retries,
	}
}

func get(destination, community, oid string, timeout time.Duration, retries int) ([]gosnmp.SnmpPDU, error) {
	conn := snmpConnection(destination, community, timeout, retries)
	err := conn.Connect()
	if err != nil {
		return nil, err
	}
	defer conn.Conn.Close()

	result, err := conn.Get([]string{oid})
	if err != nil {
		return nil, err
	}

	pdus := []gosnmp.SnmpPDU{}
	for _, pdu := range result.Variables {
		pdus = append(pdus, pdu)
	}
	return pdus, nil
}

func processQueriesFromChannel(input chan Query, processed chan Query) {
	for query := range input {
		handleQuery(&query)
		processed <- query
	}
}

type QueryWithOutputChannel struct {
	query           Query
	responseChannel chan Query
}

func ProcessAndDispatchQueries(input chan QueryWithOutputChannel, contention int) {

	inputQueries := make(chan Query, 10)
	processed := make(chan Query, 10)

	go Process(inputQueries, processed, contention)

	m := make(map[int]chan Query)
    i := 0
	for {
		select {
		case queryWithOutputChannel := <-input:
            queryWithOutputChannel.query.Id = i
            i += 1
			m[queryWithOutputChannel.query.Id] = queryWithOutputChannel.responseChannel
			inputQueries <- queryWithOutputChannel.query
		case processedQuery := <-processed:
			m[processedQuery.Id] <- processedQuery
			delete(m, processedQuery.Id)
		}
	}
}

func ExecuteQuery(queryChannel chan QueryWithOutputChannel, query Query) Query {
	output := make(chan Query)
	queryChannel <- QueryWithOutputChannel{query, output}
	processedQuery := <-output
	return processedQuery
}
