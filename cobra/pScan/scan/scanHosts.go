package scan


type state bool

type PortState struct {
	Port int
	Open state
}

func (s state) String() string {
	if s {
		return "open"
	}
	return "closed"
}

// scanPort performs a port scan on a single TCP port
func scanPort(host string, port int) PortState {
	p:=PortState {
		Port: port,
	}
}