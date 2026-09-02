package scan

import (
	"errors"
	"sort"
	"fmt"
)

var (
	ErrExists = errors.New("host already in the list")
	ErrNotExists = errors.New("host not in the list")
)

type HostsList struct{
	Hosts	[]string
}

func (hl *HostsList)search(host string)(bool, int) {
	sort.Strings(hl.Hosts)

	i:= sort.SearchStrings(hl.Hosts, host)
	if i < len(hl.Hosts) && hl.Hosts[i]==host{
		return true, i
	}
	return false, -1
}

func (hl *HostsList) Add(host string) error {
	if found, _  := hl.search(host); found {
		return fmt.Errorf("%w: %s", ErrExists, host)
	}
	hl.Hosts = append(hl.Hosts, host)
	return nil
}