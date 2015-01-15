package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/aleasoluciones/gosnmpquerier"
)

func main() {
	community := flag.String("community", "puiblic", "snmp v2 community")
	host := flag.String("host", "127.0.0.1", "host")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s: [options] [[oid] ...]\n", os.Args[0])
		flag.PrintDefaults()
	}

	flag.Parse()

	oids := []string{}
	for _, oid := range flag.Args() {
		oids = append(oids, oid)
	}

	querier := gosnmpquerier.NewSyncQuerier(1, 3, 3*time.Second)
	result, err := querier.Get(*host, *community, oids, 1*time.Second, 1)
	fmt.Println(result, err)
}
