package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
)

func signalHandler() {
	// TODO: implement save shutdown
	// ELECTRON!
	// deal with external shutdown from electron

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGHUP)
	signal.Notify(sigChan, syscall.SIGTERM)
	signal.Notify(sigChan, syscall.SIGINT)
	for sig := range sigChan {
		log.Printf("\n[server] received signal: <%v> ... shutting down", sig.String())

		// send shutdown msg to hubs
		log.Println("[server] sending shutdown message to hubs:")
		for id, hub := range wsHubs {
			log.Println(">", id)
			hub.shutdown(sig.String())
		}

		// close db
		dbClient.Close()

		os.Exit(0)
	}
}
