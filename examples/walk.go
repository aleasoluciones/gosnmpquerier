package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/aleasoluciones/gosnmpquerier"
)

func main() {
	community := flag.String("community", "alea2", "snmp v2 community")
	host := flag.String("host", "192.168.5.15", "host")
	oid := flag.String("oid", "1.3.6.1.2.1.1", "RootOid for the walk")
	flag.Parse()

	querier := gosnmpquerier.NewSyncQuerier(1, 3, 3*time.Second)
	result, err := querier.Walk(*host, *community, *oid, 1*time.Second, 1)
	if err == nil {
		for _, r := range result {
			fmt.Println(r)
		}
	}
}
