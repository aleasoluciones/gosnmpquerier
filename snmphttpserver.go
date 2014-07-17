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

func rootHandler(input chan snmpquery.QueryWithOutputChannel, w http.ResponseWriter, r *http.Request) {

    cmd, _ := snmpquery.ConvertCommand(r.FormValue("cmd"))
    query := snmpquery.Query{
        Cmd:         cmd,
        Community:   r.FormValue("community"),
        Oid:         r.FormValue("oid"),
        Timeout:     time.Duration(10) * time.Second,
        Retries:     1,
        Destination: r.FormValue("destination"),
    }
    processed := snmpquery.ExecuteQuery(input, query)
    jsonProcessed, err := snmpquery.ToJson(&processed)
    if err != nil {
        fmt.Fprint(w, err)
    }
    fmt.Fprint(w, jsonProcessed)

}

func main() {
    var input = make(chan snmpquery.QueryWithOutputChannel)
    go snmpquery.ProcessAndDispatchQueries(input, CONTENTION)

    http.HandleFunc("/", func (w http.ResponseWriter, r *http.Request) {
        rootHandler(input, w, r)
    })
    log.Fatal(http.ListenAndServe(":8080", nil))
}
