package bolt

import (
	"os"
	"testing"
)

var c *Client

func TestClient(t *testing.T) {
	c = NewClient(nil)
	c.Open()

	t.Run("PutSession", testPutSession)
	t.Run("PutUser", testPutUser)
	t.Run("GetSession", testGetSession)
	t.Run("GetUser", testGetUser)

	c.Close()
	os.Remove(c.Path)
}

func testPutSession(t *testing.T) {
	err := c.Put(SessionBucket, "test", "{\"id\":\"test\",\"data\":\"this is the test\"}")
	if err != nil {
		t.Fatal(err)
	}
}

func testPutUser(t *testing.T) {
	err := c.Put(UserBucket, "test", "{\"id\":\"test\",\"data\":\"this is the test\"}")
	if err != nil {
		t.Fatal(err)
	}
}

func testGetSession(t *testing.T) {
	_, err := c.Get(SessionBucket, "test")
	if err != nil {
		t.Fatal(err)
	}
}

func testGetUser(t *testing.T) {
	_, err := c.Get(UserBucket, "test")
	if err != nil {
		t.Fatal(err)
	}
}
