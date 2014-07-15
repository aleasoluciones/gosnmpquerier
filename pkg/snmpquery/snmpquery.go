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
	Destination string
	Response    []gosnmp.SnmpPDU
	Error       int
}

func NewQuery(id int, cmd OpSnmp, destination, community, oid string) *Query {
	return &Query{
		Id:          id,
		Cmd:         cmd,
		Community:   community,
		Oid:         oid,
		Destination: destination,
	}
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
		result, err := walk(query.Destination, query.Community, query.Oid, time.Duration(10*time.Second))
		if err == nil { // error nil means no error
			query.Response = result
		}
	case GET:
		result, err := get(query.Destination, query.Community, query.Oid, time.Duration(10*time.Second))
		if err == nil { // error nil means no error
			query.Response = result
		}
	}

}

func walk(destination, community, oid string, timeout time.Duration) ([]gosnmp.SnmpPDU, error) {
	conn := snmpConnection(destination, community, timeout)
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

func snmpConnection(destination, community string, timeout time.Duration) gosnmp.GoSNMP {
	return gosnmp.GoSNMP{
		Target:    destination,
		Port:      161,
		Community: community,
		Version:   gosnmp.Version2c,
		Timeout:   timeout,
		Retries:   1,
	}
}

func get(destination, community, oid string, timeout time.Duration) ([]gosnmp.SnmpPDU, error) {
	conn := snmpConnection(destination, community, timeout)
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
