package message

// Provider can provide raw Message data to monitor
type Provider interface {
	Emitter
	Absorber
	Verbosity
}

// Verbosity implements methods to set verbose logs
type Verbosity interface {
	SetVerbose(bool)
	GetVerbose() bool
}

// Emitter can emit []byte data
type Emitter interface {
	Emit() []byte
}

// Absorber can accept []byte data
type Absorber interface {
	Absorb([]byte)
}
