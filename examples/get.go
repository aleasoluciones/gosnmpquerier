package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/aleasoluciones/gosnmpquerier"
)

func main() {
	community := flag.String("community", "public", "snmp v2 community")
	host := flag.String("host", "127.0.0.1", "host")
	timeout := flag.Duration("timeout", 1*time.Second, "Timeout (ms/s/m/h)")
	retries := flag.Int("retries", 1, "Retries")

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
	result, err := querier.Get(*host, *community, oids, *timeout, *retries)
	fmt.Println(result, err)
}
