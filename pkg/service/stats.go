package service

import (
	"bytes"
	"fmt"
)

func StatReport() (string, error) {
	var buf bytes.Buffer

	buf.WriteString("drawr backend:\n")
	buf.WriteString("connected clients:\n")

	for id, hub := range svc.hubs {
		buf.WriteString("> " + id + ":")

		connList := hub.ListConnections()
		for _, s := range connList {
			buf.WriteString(s + "\n")
		}

		buf.WriteString("\n")
	}

	return buf.String(), nil
}

func DatabaseReport() (string, error) {
	var buf bytes.Buffer

	buf.WriteString("drawr backend:\n")
	buf.WriteString("database\n")

	buf.WriteString(fmt.Sprintf("stats: %+v\n\n", svc.db.Stats()))

	buf.WriteString(svc.db.String() + "\n")

	return buf.String(), nil
}
