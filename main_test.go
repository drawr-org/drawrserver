package main

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"./api"
	"github.com/pressly/chi"
)

// TestServer tests a running server instance
func TestServer(t *testing.T) {
	var (
		wants = "this is api version 1"
	)

	// mimick start of main()
	r := chi.NewRouter()
	api.Routing(r)

	ts := httptest.NewServer(r)
	defer ts.Close()

	res, err := http.Get(ts.URL + "/api/v1")
	if err != nil {
		t.Fatal(err)
	}

	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		t.Fatal(err)
	}

	if string(body) != wants {
		t.FailNow()
	}
}
