package api

import "github.com/drawr-team/drawrserver/pkg/bolt"

// Options holds the basic configuration of the http server
// TODO implement reading options from a config file
type Options struct {
	Port      string        `json:"port"`
	RWTimeout int64         `json:"timeout"`
	Verbose   bool          `json:"verbose"`
	Debug     bool          `json:"debug"`
	Database  *bolt.Options `json:"database"`
}
