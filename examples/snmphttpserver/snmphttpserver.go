package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/aleasoluciones/gosnmpquerier"
)

// Example of execution
// Run the http server
//		go run examples/snmphttpserver/snmphttpserver.go
// curl a request
//		curl http://127.0.0.1:8080 -X PUT  -H 'Content-Type: multipart/form-data' -d '{"cmd":"walk", "destination":"MY_HOST_IP", "community":"MY_COMMUNITY", "oid":["AN_OID"]}'

const (
	CONTENTION = 4
)

func rootHandler(querier gosnmpquerier.SyncQuerier, w http.ResponseWriter, r *http.Request) {

	cmd, _ := gosnmpquerier.ConvertCommand(r.FormValue("cmd"))
	query := gosnmpquerier.Query{
		Cmd:         cmd,
		Community:   r.FormValue("community"),
		Oids:        []string{r.FormValue("oid")},
		Timeout:     time.Duration(10) * time.Second,
		Retries:     1,
		Destination: r.FormValue("destination"),
	}
	processed := querier.ExecuteQuery(query)
	jsonProcessed, err := gosnmpquerier.ToJson(&processed)
	if err != nil {
		fmt.Fprint(w, err)
	}
	fmt.Fprint(w, jsonProcessed)

}

func main() {

	querier := gosnmpquerier.NewSyncQuerier(CONTENTION, 3, 3*time.Second)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		rootHandler(querier, w, r)
	})
	log.Fatal(http.ListenAndServe(":8080", nil))
}
