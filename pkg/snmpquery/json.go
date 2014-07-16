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

func ToJson() string {
	return `{ "body":"test" }`
}

func FromJson(jsonText string) *Query {
	var m queryMessage
	m.Timeout = 2
	m.Retries = 1

	b := []byte(jsonText)
	err := json.Unmarshal(b, &m)
	if err != nil {
		fmt.Println("Invalid jsonText format", err, jsonText)
		return nil
	}

	cmd, err := convertCommand(m.Command)
	if err != nil {
		fmt.Println("ERROR", err)
	}
	// TODO process err
	return &Query{
		Cmd:         cmd,
		Community:   m.Community,
		Oid:         m.Oid,
		Destination: m.Destination,
		Timeout:     time.Duration(m.Timeout) * time.Second,
		Retries:     m.Retries,
	}
}

func convertCommand(command string) (OpSnmp, error) {
	switch command {
	case "walk":
		return WALK, nil
	case "get":
		return GET, nil
	default:
		return 0, fmt.Errorf("Unsupported command %s ", command)
	}
}
