package message

import (
	"encoding/json"
	"log"
	"testing"

	"github.com/drawr-team/core-server/message"
	"github.com/drawr-team/core-server/test"
)

func TestNewSession(t *testing.T) {
	var p = test.MockProvider{
		mChan: make(chan []byte, 1),
		aChan: make(chan []byte, 1),
		T:     t,
	}
	var db = new(test.MockDB)

	// create a message
	mockMessage, err := message.CreateMessage(NewSessionMessageType, NewSessionData{
		Username:    "johndoe",
		SessionName: "test session",
	})
	if err != nil {
		t.Fatal(err)
	}
	// marshal it
	rawMessage, err := json.Marshal(mockMessage)
	if err != nil {
		t.Fatal(err)
	}

	if err := HandleNewSession(*m, p, db); err != nil {
		t.Fatal(err)
	}

	resp := p.Emit()
	log.Println(string(resp))
	log.Printf("%+v \n", db)
}
