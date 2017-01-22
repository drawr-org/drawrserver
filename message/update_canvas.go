package message

import (
	"encoding/json"
	"log"
	"time"

	"github.com/drawr-team/core-server/bolt"
)

const (
	UpdateCanvasMessageType       = "update-canvas"
	UpdateCanvasDataStatusSuccess = "update-canvas-success"
	UpdateCanvasDataStatusFailure = "update-canvas-failure"
	// UpdateCanvasDataMessageSuccess = ""
	// UpdateCanvasDataMessageFailure = "Session could not be joined."
)

// UpdateCanvasData is the data used to get/push canvas updates from/to the clients
type UpdateCanvasData struct {
	Username    string    `json:"username"`
	SessionID   string    `json:"sessionId"`
	CanvasState string    `json:"canvasState"`
	Timestamp   time.Time `json:"timestamp"`
}

// HandleUpdateCanvas takes the data from the client
// sets a timestamp,
// pushes the data into our database
// and pushes it on to the other clients
func HandleUpdateCanvas(m GenericMessage, p Provider, db bolt.DBClient) error {
	var data UpdateCanvasData
	if err := json.Unmarshal(m.Data, &data); err != nil {
		return err
	}

	data.Timestamp = time.Now()

	// TODO:
	// save canvas state in DB

	msg, err := NewMessage(UpdateCanvasMessageType, data)
	if err != nil {
		return err
	}

	resp, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	// log
	if p.GetVerbose() {
		// TODO where do we get the session id from?
		// this runs in session seperated hubs already
		log.Printf("[message] session: %v, user: %v\n", data.SessionID, data.Username)
	}

	p.Absorb(resp)
	return nil
}
