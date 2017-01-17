package main

import (
	"log"
	"net/http"

	"github.com/drawr-team/core-server/bolt"
)

// ValidateHandler validates the SessionID it is given
// and returns:
// * 200 OK - if the session exists in the database
// * 404 NotFound - if not
func ValidateHandler(w http.ResponseWriter, req *http.Request) {
	urlv := req.URL.Query()
	if urlv.Get("sessionId") == "" {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
	}
	// on receiving /validate?sessionId=__test__ we just return HTTP 200
	// used to test client code
	if urlv.Get("sessionId") == "__test__" {
		log.Println("[validate] testing client. Returning 200 OK")
		http.Error(w, http.StatusText(http.StatusOK), http.StatusOK)
		return
	}
	sess, err := dbClient.Get(bolt.SessionBucket, urlv.Get("sessionId"))
	if err != nil {
		if err == bolt.ErrNotFound {
			// if session id is not in database
			log.Println("[validate] session id not found:", urlv.Get("sessionId"))
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
		// if bolt returns an error we don't expect
		log.Println("[validate] unknown error", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	// id everything is fine
	log.Println("[validate] session id exists:", urlv.Get("sessionId"))
	w.Write([]byte("valid session id:\n"))
	w.Write(sess)
	http.Error(w, http.StatusText(http.StatusOK), http.StatusOK)
	return
}
