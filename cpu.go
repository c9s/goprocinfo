package goproc

import (
	"io/ioutil"
	"strconv"
	"strings"
)

type Stat struct {
	CPUStatAll *CPUStat
	CPUStats   []*CPUStat
}

type CPUStat struct {
	Id        string
	User      uint64
	Nice      uint64
	System    uint64
	Idle      uint64
	IOWait    uint64
	IRQ       uint64
	SoftIRQ   uint64
	Steal     uint64
	Guest     uint64
	GuestNice uint64
}

func parseCPUStat(line string) *CPUStat {
	fields := strings.Fields(line)
	s := CPUStat{}
	s.User, _ = strconv.ParseUint(fields[1], 10, 32)
	s.Nice, _ = strconv.ParseUint(fields[2], 10, 32)
	s.System, _ = strconv.ParseUint(fields[3], 10, 32)
	s.Idle, _ = strconv.ParseUint(fields[4], 10, 32)
	s.IOWait, _ = strconv.ParseUint(fields[5], 10, 32)
	s.IRQ, _ = strconv.ParseUint(fields[6], 10, 32)
	s.SoftIRQ, _ = strconv.ParseUint(fields[7], 10, 32)
	s.Steal, _ = strconv.ParseUint(fields[8], 10, 32)
	s.Guest, _ = strconv.ParseUint(fields[9], 10, 32)
	s.GuestNice, _ = strconv.ParseUint(fields[10], 10, 32)
	return &s
}

func ReadStat(path string) (*Stat, error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	content := string(b)
	lines := strings.Split(content, "\n")

	var stat Stat = Stat{}

	for _, line := range lines {
		if strings.HasPrefix(line, "cpu") {
			if cpuStat := parseCPUStat(line); cpuStat != nil {
				if stat.CPUStatAll == nil {
					stat.CPUStatAll = cpuStat
				} else {
					stat.CPUStats = append(stat.CPUStats, cpuStat)
				}
			}
		}
	}
	return &stat, nil
}
