package linux

import (
	"fmt"
	"io/ioutil"
	"strings"
)

type MemInfo struct {
	Total      uint64
	Free       uint64
	Buffers    uint64
	Cached     uint64
	SwapCached uint64
	Active     uint64
	Inactive   uint64
}

func ReadMemInfo(path string) (*MemInfo, error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	lines := strings.Split(string(b), "\n")
	meminfo := MemInfo{}
	if _, err := fmt.Sscanf(lines[0], "MemTotal: %d kB", &meminfo.Total); err != nil {
		return nil, err
	}
	if _, err := fmt.Sscanf(lines[1], "MemFree: %d kB", &meminfo.Free); err != nil {
		return nil, err
	}
	if _, err := fmt.Sscanf(lines[2], "Buffers: %d kB", &meminfo.Buffers); err != nil {
		return nil, err
	}
	if _, err := fmt.Sscanf(lines[3], "Cached: %d kB", &meminfo.Cached); err != nil {
		return nil, err
	}
	if _, err := fmt.Sscanf(lines[4], "SwapCached: %d kB", &meminfo.SwapCached); err != nil {
		return nil, err
	}
	if _, err := fmt.Sscanf(lines[5], "Active: %d kB", &meminfo.Active); err != nil {
		return nil, err
	}
	if _, err := fmt.Sscanf(lines[6], "Inactive: %d kB", &meminfo.Inactive); err != nil {
		return nil, err
	}
	return &meminfo, nil
}
