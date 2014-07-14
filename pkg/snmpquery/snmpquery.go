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

func printValue(pdu gosnmp.SnmpPDU) error {
	fmt.Printf("%s = ", pdu.Name)

	switch pdu.Type {
	case gosnmp.OctetString:
		fmt.Printf("STRING: %s\n", pdu.Value.(string))
	case gosnmp.Counter64:
		fmt.Printf("COUNTER64 %d: %d\n", pdu.Type, gosnmp.ToBigInt(pdu.Value))
	default:
		fmt.Printf("TYPE %d: %d\n", pdu.Type, gosnmp.ToBigInt(pdu.Value))
	}
	return nil
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

	err = gosnmp.Default.BulkWalk(query.Oid, printValue)
	if err != nil {
		fmt.Printf("Walk Error: %v\n", err)
		return
	}

	query.Response = "whatever " + strconv.Itoa(rand.Intn(1e3))
}

func processQueriesFromChannel(input chan Query, processed chan Query) {
	for query := range input {
		handleQuery(&query)
		processed <- query
	}
}
