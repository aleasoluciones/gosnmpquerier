package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/aleasoluciones/gosnmpquerier"
)

func main() {
	community := flag.String("community", "public", "snmp v2 community")
	host := flag.String("host", "127.0.0.1", "host")
	oid := flag.String("oid", "1.3.6.1.2.1.1.1.0", "Oid to get")
	flag.Parse()

	querier := gosnmpquerier.NewSyncQuerier(1)
	result, err := querier.Get(*host, *community, []string{*oid}, 1*time.Second, 1)
	result2, err2 := querier.Get(*host, *community, []string{*oid}, 1*time.Second, 1)
	fmt.Println(result, err)
	fmt.Println()
	fmt.Println(result2, err2)
}
