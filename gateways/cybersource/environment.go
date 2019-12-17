package cybersource

type Environment struct {
	host string
}

func NewEnvironment(host string) Environment {
	return Environment{host: host}
}

func (e Environment) Host() string {
	return e.host
}
