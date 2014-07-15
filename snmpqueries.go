package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/eferro/go-snmpqueries/pkg/snmpquery"
)

const (
	CONTENTION = 4
)

func readQueriesFromStdin(input chan snmpquery.Query) {
	reader := bufio.NewReader(os.Stdin)
	queryId := 0
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}

		fields := strings.Fields(line)
		if len(fields) != 4 {
			log.Println("InvalidLine", line)
		} else {
			var cmd snmpquery.OpSnmp
			switch fields[0] {
			case "walk":
				cmd = snmpquery.WALK
			case "get":
				cmd = snmpquery.GET
			default:
				log.Println("InvalidLine", line)
				continue
			}

			timeout := time.Duration(2 * time.Second)
			retries := 2

			query := snmpquery.Query{
				Id:          queryId,
				Cmd:         cmd,
				Community:   fields[2],
				Oid:         fields[3],
				Destination: fields[1],
				Timeout:     timeout,
				Retries:     retries,
			}
			input <- query
			queryId += 1
		}
	}
}

func printResults(processed chan snmpquery.Query) {
	for query := range processed {
		fmt.Println("Result", query)
	}
}

func main() {
	input := make(chan snmpquery.Query, 10)
	processed := make(chan snmpquery.Query, 10)

	go readQueriesFromStdin(input)
	go snmpquery.Process(input, processed, CONTENTION)

	printResults(processed)
}
