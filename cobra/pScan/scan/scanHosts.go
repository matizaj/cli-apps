package scan

import (
	"fmt"
	"net"
	"time"
)

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

type Result struct {
	Host       string
	NotFound   bool
	PortStates []PortState
}

// scanPort performs a port scan on a single TCP port
func scanPort(host string, port int) PortState {
	p := PortState{
		Port: port,
	}
	address := net.JoinHostPort(host, fmt.Sprintf("%d", port))
	c, err := net.DialTimeout("tcp", address, 1*time.Second)
	if err != nil {
		return p
	}
	c.Close()
	p.Open = true
	return p
}

func Run(hl *HostsList, ports []int) []Result {
	res := make([]Result, 0, len(hl.Hosts))

	for _, h := range hl.Hosts {
		r := Result{
			Host: h,
		}
		_, err := net.LookupHost(h)
		if err != nil {
			r.NotFound = true
			res = append(res, r)
			continue
		}

		for _, p := range ports {
			r.PortStates = append(r.PortStates, scanPort(h, p))
		}

		res = append(res, r)
	}
	return res
}
