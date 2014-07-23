package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/eferro/gosnmpquerier"
)

const (
	CONTENTION = 4
)

func rootHandler(querier *gosnmpquerier.SyncQuerier, w http.ResponseWriter, r *http.Request) {

	cmd, _ := gosnmpquerier.ConvertCommand(r.FormValue("cmd"))
	query := gosnmpquerier.Query{
		Cmd:         cmd,
		Community:   r.FormValue("community"),
		Oid:         r.FormValue("oid"),
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

	querier := gosnmpquerier.NewSyncQuerier(CONTENTION)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		rootHandler(querier, w, r)
	})
	log.Fatal(http.ListenAndServe(":8080", nil))
}
