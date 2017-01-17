package message

// Provider can provide raw Message data to monitor
type Provider interface {
	Emitter
	Absorber
}

// Emitter can emit []byte data
type Emitter interface {
	Emit() []byte
}

// Absorber can accept []byte data
type Absorber interface {
	Absorb([]byte)
	AbsorbTo(string, []byte)
}
