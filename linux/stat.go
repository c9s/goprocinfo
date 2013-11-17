package linux

import (
	"io/ioutil"
	"strconv"
	"strings"
	"time"
)

type Stat struct {
	CPUStatAll      CPUStat
	CPUStats        []CPUStat
	Interrupts      uint64
	ContextSwitches uint64
	BootTime        time.Time
	Processes       uint64
	ProcsRunning    uint64
	ProcsBlocked    uint64
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

func createCPUStat(fields []string) *CPUStat {
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

	for i, line := range lines {
		fields := strings.Fields(line)
		if len(fields) == 0 {
			continue
		}
		if fields[0] == "cpu" {
			if cpuStat := createCPUStat(fields); cpuStat != nil {
				if i == 0 {
					stat.CPUStatAll = *cpuStat
				} else {
					stat.CPUStats = append(stat.CPUStats, *cpuStat)
				}
			}
		} else if fields[0] == "intr" {
			stat.Interrupts, _ = strconv.ParseUint(fields[1], 10, 64)
		} else if fields[0] == "ctxt" {
			stat.ContextSwitches, _ = strconv.ParseUint(fields[1], 10, 64)
		} else if fields[0] == "btime" {
			seconds, _ := strconv.ParseInt(fields[1], 10, 64)
			stat.BootTime = time.Unix(seconds, 0)
		} else if fields[0] == "processes" {
			stat.Processes, _ = strconv.ParseUint(fields[1], 10, 64)
		} else if fields[0] == "procs_running" {
			stat.ProcsRunning, _ = strconv.ParseUint(fields[1], 10, 64)
		} else if fields[0] == "procs_blocked" {
			stat.ProcsBlocked, _ = strconv.ParseUint(fields[1], 10, 64)
		}
	}
	return &stat, nil
}
