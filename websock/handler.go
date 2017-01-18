package websock

import (
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

// Handler implements the http.Handler interface
// for a websock.Hub
type Handler struct {
	Hub *Hub
}

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin:     func(*http.Request) bool { return true },
	}

	// get the websocket connection
	wsConn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		if h.Hub.Verbose {
			log.Println("error upgrading:", err)
		}
		// return NotFound status if upgrading fails
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	// TODO: review this. there might be a better way to do this

	c := NewConnection(h.Hub, wsConn)
	c.Hub.AddConnection(c)
	defer c.Hub.RemoveConnection(c)

	var wg sync.WaitGroup
	wg.Add(2)

	go c.Writer(&wg, wsConn)
	go c.Reader(&wg, wsConn)

	wg.Wait()
	wsConn.Close()
}
