package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/drawr-team/drawrserver/cmd/drawrserver/api"

	log "github.com/golang/glog"
)

var (
	port         = flag.String("port", "3000", "port to run the server on")
	dbPath       = flag.String("database", "data.db", "location of the database file")
	printversion = flag.Bool("version", false, "print version number")
	configFile   = flag.String("config", "config.toml", "path to the config file")
)

func init() {
	flag.Parse()
	if *printversion {
		fmt.Print(version)
		os.Exit(0)
	}
}

func main() {
	config, err := loadConfig(*configFile)
	if err != nil {
		log.Error("Error loading configuration: ", err)
	}
	log.Info("Config loaded")
	log.V(2).Infof("%+v", config)

	var server http.Server
	server.Addr = fmt.Sprintf("%s:%d", config.Server.Host, config.Server.Port)
	server.ReadTimeout = time.Duration(config.Server.RWTimeout) * time.Second
	server.WriteTimeout = time.Duration(config.Server.RWTimeout) * time.Second

	// init api pkg
	handler, err := api.Init(config.Database.Path, config.Database.Timeout)
	if err != nil {
		log.Fatal("Error setting up API: ", err)
	}
	server.Handler = handler
	log.Info("API initialized")

	// register signals
	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, syscall.SIGHUP)
	signal.Notify(sigChan, syscall.SIGTERM)
	signal.Notify(sigChan, syscall.SIGINT)

	// start server
	go func() {
		if err := server.ListenAndServe(); err != nil {
			log.Fatal("Error serving content: ", err)
		}
	}()
	log.Info("Server started at ", server.Addr)

	// wait on signal channel
	<-sigChan
	log.Warning("Server shutting down")

	ctx, cancelFunc := context.WithTimeout(context.Background(), 5*time.Second)
	cancelFunc()
	server.Shutdown(ctx)
	log.Info("Server stopped")
}
