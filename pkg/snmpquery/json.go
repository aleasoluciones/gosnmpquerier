package snmpquery

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/soniah/gosnmp"
)

type queryMessage struct {
	Command        string
	Destination    string
	Community      string
	Oid            string
	Timeout        int
	Retries        int
	AdditionalInfo interface{}
}

type outputMessage struct {
	Id          int
	Command     OpSnmp
	Community   string
	Oid         string
	Timeout     time.Duration
	Retries     int
	Destination string
	Response    []gosnmp.SnmpPDU
	Error       string
}

func ToJson(query *Query) (string, error) {

    d := outputMessage{
        Id: query.Id,
        Command: query.Cmd,
        Community: query.Community,
        Oid: query.Oid,
        Timeout: query.Timeout,
        Retries: query.Retries,
        Destination: query.Destination,
        Response: query.Response,
        Error: query.Error.Error(),
    }
    fmt.Println(d)
	b, err := json.Marshal(d)

	if err != nil {
		return "", err
	}

	return string(b), nil
}

func FromJson(jsonText string) (*Query, error) {
	var m queryMessage
	m.Timeout = 2
	m.Retries = 1

	b := []byte(jsonText)
	err := json.Unmarshal(b, &m)
	if err != nil {
		return nil, err
	}

	cmd, err := ConvertCommand(m.Command)
	if err != nil {
		return nil, err
	}

	q := Query{
		Cmd:         cmd,
		Community:   m.Community,
		Oid:         m.Oid,
		Destination: m.Destination,
		Timeout:     time.Duration(m.Timeout) * time.Second,
		Retries:     m.Retries,
	}
	return &q, nil
}

func ConvertCommand(command string) (OpSnmp, error) {
	switch command {
	case "walk":
		return WALK, nil
	case "get":
		return GET, nil
	default:
		return 0, fmt.Errorf("Unsupported command %s ", command)
	}
}
