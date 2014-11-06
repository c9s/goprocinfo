package linux

import (
	"io/ioutil"
	"strconv"
	"strings"
	"time"
)

type Stat struct {
	CPUStatAll      CPUStat   `json:"cpu_all"`
	CPUStats        []CPUStat `json:"cpus"`
	Interrupts      uint64    `json:"intr"`
	ContextSwitches uint64    `json:"ctxt"`
	BootTime        time.Time `json:"btime"`
	Processes       uint64    `json:"processes"`
	ProcsRunning    uint64    `json:"procs_running"`
	ProcsBlocked    uint64    `json:"procs_blocked"`
}

type CPUStat struct {
	Id        string `json:"id"`
	User      uint64 `json:"user"`
	Nice      uint64 `json:"nice"`
	System    uint64 `json:"system"`
	Idle      uint64 `json:"idle"`
	IOWait    uint64 `json:"iowait"`
	IRQ       uint64 `json:"irq"`
	SoftIRQ   uint64 `json:"softirq"`
	Steal     uint64 `json:"steal"`
	Guest     uint64 `json:"guest"`
	GuestNice uint64 `json:"guest_nice"`
}

func createCPUStat(fields []string) *CPUStat {
	s := CPUStat{}
	s.User, _ = strconv.ParseUint(fields[1], 10, 64)
	s.Nice, _ = strconv.ParseUint(fields[2], 10, 64)
	s.System, _ = strconv.ParseUint(fields[3], 10, 64)
	s.Idle, _ = strconv.ParseUint(fields[4], 10, 64)
	s.IOWait, _ = strconv.ParseUint(fields[5], 10, 64)
	s.IRQ, _ = strconv.ParseUint(fields[6], 10, 64)
	s.SoftIRQ, _ = strconv.ParseUint(fields[7], 10, 64)
	s.Steal, _ = strconv.ParseUint(fields[8], 10, 64)
	s.Guest, _ = strconv.ParseUint(fields[9], 10, 64)
	if len(fields) > 10 {
		s.GuestNice, _ = strconv.ParseUint(fields[10], 10, 64)
	}
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
		if fields[0][:3] == "cpu" {
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
