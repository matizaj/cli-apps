package scan

import (
	"errors"
	"fmt"
	"os"
	"sort"
	"bufio"
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

func (hl *HostsList) Remove(host string) error {
	found, idx  := hl.search(host);
	if !found {
		return fmt.Errorf("%w: %s", ErrNotExists, host)
	}
	hl.Hosts = append(hl.Hosts[:idx], hl.Hosts[idx+1:]...)
	return nil
}

func (hl *HostsList) Load(hostFile string) error {
	f, err := os.Open(hostFile)
	if err != nil {
		if errors.Is(err, ErrNotExists) {
			return nil
		}
		return err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		hl.Hosts = append(hl.Hosts,scanner.Text())
	}
	return nil
}

func (hl *HostsList) Save(hostFile string) error {
	output:= ""
	for _, h := range hl.Hosts {
		output+=fmt.Sprintln(h)
	}

	return os.WriteFile(hostFile,[]byte(output), 0644)
}