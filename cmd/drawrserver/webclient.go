package main

import "net/http"

func WebClientHandler(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte("come back later to have a full blown drawr web client here..."))
	return
}
