package snmpquery

import (
	"encoding/json"
	"fmt"
	"time"
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

func ToJson(query *Query) (string, error) {

	b, err := json.Marshal(query)
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
