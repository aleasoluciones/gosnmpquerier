package snmpquery

import (
	"fmt"
	"math/rand"
	"strconv"
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
	Response    string
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

	gosnmp.Default.Community = query.Community
	gosnmp.Default.Target = query.Destination
	gosnmp.Default.Timeout = time.Duration(10 * time.Second)
	err := gosnmp.Default.Connect()
	if err != nil {
		fmt.Printf("Connect err: %v\n", err)
		return
	}
	defer gosnmp.Default.Conn.Close()

	output := make(chan string, 10)

	fn := func(pdu gosnmp.SnmpPDU) error {
		out := pdu.Name

		switch pdu.Type {
		case gosnmp.OctetString:
			out += "STRING "
		case gosnmp.Counter64:
			out += "COUNTER64 "
		}
		output <- out
		return nil
	}

	err = gosnmp.Default.BulkWalk(query.Oid, fn)
	if err != nil {
		fmt.Printf("Walk Error: %v\n", err)
		return
	}

	for result := range output {
		fmt.Println("EFA", result)
	}

	query.Response = "whatever " + strconv.Itoa(rand.Intn(1e3))
}

func processQueriesFromChannel(input chan Query, processed chan Query) {
	for query := range input {
		handleQuery(&query)
		processed <- query
	}
}
