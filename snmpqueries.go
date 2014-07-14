package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

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

			query := snmpquery.NewQuery(queryId, cmd, fields[1], fields[2], fields[3])
			input <- *query
			queryId += 1
		}
	}
}

func printResults(processed chan snmpquery.Query) {
	for query := range processed {
		fmt.Println(query)
	}
}

func main() {
	input := make(chan snmpquery.Query, 10)
	processed := make(chan snmpquery.Query, 10)

	go readQueriesFromStdin(input)
	go snmpquery.Process(input, processed, CONTENTION)

	printResults(processed)
}
