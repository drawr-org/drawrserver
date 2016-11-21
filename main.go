package main

import (
	"log"

	"./drawr"
)

func main() {
	s := drawr.NewServer()

	log.Println("Listening on...", s.Addr)
	panic(s.ListenAndServe())
}
