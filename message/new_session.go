package message

import (
	"encoding/json"
	"log"

	"github.com/drawr-team/core-server/bolt"
	"github.com/drawr-team/core-server/ulidgen"
)

const (
	NewSessionMessageType       = "new-session"
	NewSessionAckType           = "ack-session"
	NewSessionDataStatusSuccess = "new-session-success"
	NewSessionDataStatusFailure = "new-session-failure"
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
func HandleNewSession(m GenericMessage, p Provider, db bolt.DBClient) error {
	var data NewSessionData
	if err := json.Unmarshal(m.Data, &data); err != nil {
		return err
	}

	// generate new session id
	sessionID := ulidgen.Now().String()

	// create database entry for new session
	db.Put(bolt.SessionBucket, sessionID, m.Data)

	// message to clients
	var resp *GenericMessage
	if err := db.Put(bolt.SessionBucket, sessionID, data); err != nil {
		if err == bolt.ErrExists {
			// create fail response
			if resp, err = CreateMessage(NewSessionAckType, GenericAck{
				Status: NewSessionDataStatusFailure,
				Data:   NewSessionAckData{SessionID: sessionID},
			}); err != nil {
				return err
			}
		}
		log.Println(err)
	} else {
		// create success repsonse
		if resp, err = CreateMessage(NewSessionAckType, GenericAck{
			Status: NewSessionDataStatusSuccess,
			Data:   NewSessionAckData{SessionID: sessionID},
		}); err != nil {
			return err
		}
	}

	message, err := json.Marshal(resp)
	if err != nil {
		return err
	}

	p.AbsorbTo(sessionID, message)

	return nil
}
