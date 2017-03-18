package bolt

// Options holds the basic database configuration
type Options struct {
	Path    string `json:"path"`
	Timeout int64  `json:"timeout"`
	Verbose bool   `json:"verbose"`
}
