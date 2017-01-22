package main

import "log"

const version = "0.2.1"

func main() {
	defer dbClient.Close()
	go signalHandler()

	log.Println("Listening on...", server.Addr)
	if err := server.ListenAndServe(); err != nil {
		panic(err)
	}
}
