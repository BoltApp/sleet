package cybersource

// Environment specifies the host (test/production) used for API calls
type Environment struct {
	host string
}

// NewEnvironment constructs an environment based on host
func NewEnvironment(host string) Environment {
	return Environment{host: host}
}

// Host returns host specified by Environment
func (e Environment) Host() string {
	return e.host
}
