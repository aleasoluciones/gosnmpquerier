// Copyright 2014 The GoSNMPQuerier Authors. All rights reserved.  Use of this
// source code is governed by a MIT-style license that can be found in the
// LICENSE file.

package gosnmpquerier

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
    Oids           []string
    Timeout        int
    Retries        int
    AdditionalInfo interface{}
}

type outputMessage struct {
    Id          int
    Command     string
    Community   string
    Oids        []string
    Timeout     time.Duration
    Retries     int
    Destination string
    Response    []gosnmp.SnmpPDU
    Error       string
}

func ToJson(query *Query) (string, error) {
    var errString string = ""
    if query.Error != nil {
        errString = query.Error.Error()
    }
    d := outputMessage{
        Id:          query.Id,
        Command:     convertCommandToCommandString(query.Cmd),
        Community:   query.Community,
        Oids:        query.Oids,
        Timeout:     query.Timeout,
        Retries:     query.Retries,
        Destination: query.Destination,
        Response:    query.Response,
        Error:       errString,
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
        Oids:        m.Oids,
        Destination: m.Destination,
        Timeout:     time.Duration(m.Timeout) * time.Second,
        Retries:     m.Retries,
    }
    return &q, nil
}

func convertCommandToCommandString(command OpSnmp) string {
    switch command {
    case WALK:
        return "walk"
    case GET:
        return "get"
    }
    return ""
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
