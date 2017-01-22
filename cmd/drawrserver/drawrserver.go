package main

import "log"

func main() {
	defer dbClient.Close()
	go signalHandler()

	log.Println("Listening on...", server.Addr)
	if err := server.ListenAndServe(); err != nil {
		panic(err)
	}
}
