package main

import (
    "fmt"
    "log"
    "net/http"
    "time"

	"github.com/eferro/go-snmpqueries/pkg/snmpquery"
)

func rootHandler(w http.ResponseWriter, r *http.Request) {

    fmt.Println(r.FormValue("query"))

    query := snmpquery.Query{
        Cmd:         snmpquery.ConvertCommand(r.FormValue("cmd")),
        Community:   r.FormValue("community"),
        Oid:         r.FormValue("oid"),
        Timeout:     time.Duration(m.Timeout) * time.Second,
        Retries:     1,
        Destination: r.FormValue("destination"),
    }

    processed := snmpquery.ExecuteQuery(query)

    if err := encoder.Encode(&query); err != nil {
            http.Error(w, fmt.Sprintf("Cannot encode response data: %v", err), 500)
        }
    }

func main() {
    http.HandleFunc("/", rootHandler)
    log.Fatal(http.ListenAndServe(":8080", nil))
}
