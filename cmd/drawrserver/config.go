package main

import (
	"os"

	"github.com/BurntSushi/toml"
	log "github.com/golang/glog"
)

type Config struct {
	Server   ServerOptions
	Database DatabaseOptions
}

type ServerOptions struct {
	Host      string
	Port      int
	RWTimeout int
}

type DatabaseOptions struct {
	Path    string
	Timeout int
}

var defaultConfig = Config{
	Server:   defaultServerOptions,
	Database: defaultDatabaseOptions,
}
var defaultServerOptions = ServerOptions{
	Host:      "localhost",
	Port:      3000,
	RWTimeout: 5,
}
var defaultDatabaseOptions = DatabaseOptions{
	Path:    "data.db",
	Timeout: 1,
}

func loadConfig(fpath string) (*Config, error) {
	if _, err := os.Stat(fpath); os.IsNotExist(err) {
		log.Warning("No config file. Using defaults")
		return &defaultConfig, nil
	}

	var c Config
	if _, err := toml.DecodeFile(fpath, &c); err != nil {
		return nil, err
	}
	return &c, nil
}
