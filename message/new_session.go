package message

import (
	"encoding/json"
	"log"

	"github.com/drawr-team/core-server/bolt"
	"github.com/drawr-team/core-server/ulidgen"
)

const (
	NewSessionMessageType        = "new-session"
	NewSessionAckType            = "ack-session"
	NewSessionDataStatusSuccess  = "new-session-success"
	NewSessionDataStatusFailure  = "new-session-failure"
	NewSessionDataMessageSuccess = "Session successfully created."
	NewSessionDataMessageFailure = "Session could not be created."
)

// NewSessionData is the data used to initialize a new Session
type NewSessionData struct {
	Username    string `json:"username"`
	SessionName string `json:"sessionName"`
}

// NewSessionAckData is the payload of the Ack
type NewSessionAckData struct {
	SessionID string `json:"sessionId"`
}

// HandleNewSession handles a `new-session` type
func HandleNewSession(m GenericMessage, p Provider) error {
	var data NewSessionData
	if err := json.Unmarshal(m.Data, &data); err != nil {
		return err
	}

	sessionID := ulidgen.Now().String()

	if err := bolt.Put(bolt.SessionBucket, sessionID, data); err != nil {
		log.Println(err)
	}
	// TODO: session logic here

	// create repsonse
	resp, err := CreateMessage(NewSessionAckType, GenericAck{
		Status:  NewSessionDataStatusSuccess,
		Message: NewSessionDataMessageSuccess,
		Data:    NewSessionAckData{SessionID: sessionID},
	})
	if err != nil {
		return err
	}

	message, err := json.Marshal(resp)
	if err != nil {
		return err
	}

	p.Absorb(message)

	return nil
}
