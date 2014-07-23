package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/eferro/go-snmpqueries/pkg/snmpquery"
)

const (
	CONTENTION = 4
)

func rootHandler(querier *snmpquery.SyncQuerier, w http.ResponseWriter, r *http.Request) {

	cmd, _ := snmpquery.ConvertCommand(r.FormValue("cmd"))
	query := snmpquery.Query{
		Cmd:         cmd,
		Community:   r.FormValue("community"),
		Oid:         r.FormValue("oid"),
		Timeout:     time.Duration(10) * time.Second,
		Retries:     1,
		Destination: r.FormValue("destination"),
	}
	processed := querier.ExecuteQuery(query)
	jsonProcessed, err := snmpquery.ToJson(&processed)
	if err != nil {
		fmt.Fprint(w, err)
	}
	fmt.Fprint(w, jsonProcessed)

}

func main() {

	querier := snmpquery.NewSyncQuerier(CONTENTION)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		rootHandler(querier, w, r)
	})
	log.Fatal(http.ListenAndServe(":8080", nil))
}
