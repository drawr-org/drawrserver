package websock

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/websocket"
)

var testHub *Hub

type wsHandler struct{ *testing.T }
type wsServer struct {
	*httptest.Server
	URL string
}

var wsDialer = websocket.Dialer{
	Subprotocols:    []string{"drawrv1"},
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

const wsRequestURI = "/a/b?x=y"

func newTestServer(t *testing.T) *wsServer {
	testHub = NewHub()

	var s wsServer
	s.Server = httptest.NewServer(wsHandler{t})
	s.Server.URL += wsRequestURI
	s.URL = "ws" + strings.TrimPrefix(s.Server.URL, "http")
	return &s
}

func (t wsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.RequestURI() != wsRequestURI {
		t.Logf("path=%v, want %v", r.URL.RequestURI(), wsRequestURI)
		http.Error(w, "bad path", 400)
		return
	}
	c, err := Upgrade(w, r, w.Header())
	if err != nil {
		t.Log("upgrading connection failed:", err)
		http.Error(w, "upgrader failed", 500)
	}
	testHub.AddConnection(fmt.Sprintf("conn_%d", len(testHub.connections)), *c)
}

func newTestConn(t *testing.T, ts *wsServer) Connection {
	cli := websocket.Dialer{}
	_, resp, err := cli.Dial(ts.URL, nil)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != 101 {
		t.Logf("status=%v, want 101", resp.StatusCode)
		t.FailNow()
	}

	return testHub.connections[fmt.Sprintf("conn_%d", len(testHub.connections)-1)]
}

func TestUpgrader(t *testing.T) {
	ts := newTestServer(t)
	defer ts.Close()

	c := newTestConn(t, ts)
	if c.messenger == nil {
		t.Log("connection was nil")
		t.FailNow()
	}
}
