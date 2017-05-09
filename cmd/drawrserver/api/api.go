// Package api implements the HTTP API for the drawrserver
package api

import (
	"net/http"

	"github.com/drawr-team/drawrserver/pkg/service"

	log "github.com/golang/glog"
)

const version string = "v1"

// Init sets up the package
func Init(dbPath string, dbTimeout int) (http.Handler, error) {
	service.Init(dbPath, dbTimeout)

	routes, err := setupRoutes()
	if err != nil {
		log.Error("Error setting up router:", err)
		return nil, err
	}

	return routes, nil
}
