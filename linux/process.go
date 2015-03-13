package linux

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
)

type Process struct {
	Status ProcessStatus `json:"status"`
	Statm  ProcessStatm  `json:"statm"`
	Stat   ProcessStat   `json:"stat"`
	IO     ProcessIO     `json:"io"`
}

type ProcessStatus struct{}

// I/O statistics for the process.
type ProcessIO struct {
	RChar               uint64 `json:"rchar" field:"rchar"`                                 // chars read
	WChar               uint64 `json:"wchar" field:"wchar"`                                 // chars written
	Syscr               uint64 `json:"syscr" field:"syscr"`                                 // read syscalls
	Syscw               uint64 `json:"syscw" field:"syscw"`                                 // write syscalls
	ReadBytes           uint64 `json:"read_bytes" field:"read_bytes"`                       // bytes read
	WriteBytes          uint64 `json:"write_bytes" field:"write_bytes"`                     // bytes written
	CancelledWriteBytes uint64 `json:"cancelled_write_bytes" field:"cancelled_write_bytes"` // bytes truncated
}

// Provides information about memory usage, measured in pages.
type ProcessStatm struct {
	Size     uint64 `json:"size"`     // total program size
	Resident uint64 `json:"resident"` // resident set size
	Share    uint64 `json:"share"`    // shared pages
	Text     uint64 `json:"text"`     // text (code)
	Lib      uint64 `json:"lib"`      // library (unused in Linux 2.6)
	Data     uint64 `json:"data"`     // data + stack
	Dirty    uint64 `json:"dirty"`    // dirty pages (unused in Linux 2.6)
}

// Status information about the process.
type ProcessStat struct {
	Pid                 uint64 `json:"pid"`
	Comm                string `json:"comm"`
	State               string `json:"state"`
	Ppid                int64  `json:"ppid"`
	Pgrp                int64  `json:"pgrp"`
	Session             int64  `json:"session"`
	TtyNr               int64  `json:"tty_nr"`
	Tpgid               int64  `json:"tpgid"`
	Flags               uint64 `json:"flags"`
	Minflt              uint64 `json:"minflt"`
	Cminflt             uint64 `json:"cminflt"`
	Majflt              uint64 `json:"majflt"`
	Cmajflt             uint64 `json:"cmajflt"`
	Utime               uint64 `json:"utime"`
	Stime               uint64 `json:"stime"`
	Cutime              int64  `json:"cutime"`
	Cstime              int64  `json:"cstime"`
	Priority            int64  `json:"priority"`
	Nice                int64  `json:"nice"`
	NumThreads          int64  `json:"num_threads"`
	Itrealvalue         int64  `json:"itrealvalue"`
	Starttime           uint64 `json:"starttime"`
	Vsize               uint64 `json:"vsize"`
	Rss                 int64  `json:"rss"`
	Rsslim              uint64 `json:"rsslim"`
	Startcode           uint64 `json:"startcode"`
	Endcode             uint64 `json:"endcode"`
	Startstack          uint64 `json:"startstack"`
	Kstkesp             uint64 `json:"kstkesp"`
	Kstkeip             uint64 `json:"kstkeip"`
	Signal              uint64 `json:"signal"`
	Blocked             uint64 `json:"blocked"`
	Sigignore           uint64 `json:"sigignore"`
	Sigcatch            uint64 `json:"sigcatch"`
	Wchan               uint64 `json:"wchan"`
	Nswap               uint64 `json:"nswap"`
	Cnswap              uint64 `json:"cnswap"`
	ExitSignal          int64  `json:"exit_signal"`
	Processor           int64  `json:"processor"`
	RtPriority          uint64 `json:"rt_priority"`
	Policy              uint64 `json:"policy"`
	DelayacctBlkioTicks uint64 `json:"delayacct_blkio_ticks"`
	GuestTime           uint64 `json:"guest_time"`
	CguestTime          int64  `json:"cguest_time"`
	StartData           uint64 `json:"start_data"`
	EndData             uint64 `json:"end_data"`
	StartBrk            uint64 `json:"start_brk"`
	ArgStart            uint64 `json:"arg_start"`
	ArgEnd              uint64 `json:"arg_end"`
	EnvStart            uint64 `json:"env_start"`
	EnvEnd              uint64 `json:"env_end"`
	ExitCode            int64  `json:"exit_code"`
}

func ReadProcess(pid uint64, path string) (*Process, error) {

	p := filepath.Join(path, strconv.FormatUint(pid, 10))

	if _, err := os.Stat(p); err != nil {
		return nil, err
	}

	process := Process{}

	statusPath := filepath.Join(p, "status")

	var err error
	var io *ProcessIO
	var statm *ProcessStatm
	var stat *ProcessStat

	if io, err = ReadProcessIO(filepath.Join(p, "io")); err != nil {
		return nil, err
	}

	if statm, err = ReadProcessStatm(filepath.Join(p, "statm")); err != nil {
		return nil, err
	}

	if stat, err = ReadProcessStat(filepath.Join(p, "stat")); err != nil {
		return nil, err
	}

	_ = statusPath

	process.IO = *io
	process.Statm = *statm
	process.Stat = *stat

	return &process, nil
}

func ReadProcessStatus(path string) (*ProcessStatus, error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	status := ProcessStatus{}
	_ = b
	return &status, nil
}

func ReadProcessIO(path string) (*ProcessIO, error) {

	b, err := ioutil.ReadFile(path)

	if err != nil {
		return nil, err
	}

	// Maps a io metric to its value (i.e. rchar --> 100000)
	m := map[string]uint64{}

	var io ProcessIO = ProcessIO{}

	lines := strings.Split(string(b), "\n")

	for _, line := range lines {

		if strings.Index(line, ": ") == -1 {
			continue
		}

		l := strings.Split(line, ": ")

		k := l[0]
		v, err := strconv.ParseUint(l[1], 10, 64)

		if err != nil {
			return nil, err
		}

		m[k] = v

	}

	e := reflect.ValueOf(&io).Elem()
	t := e.Type()

	for i := 0; i < e.NumField(); i++ {

		k := t.Field(i).Tag.Get("field")

		v, ok := m[k]

		if ok {
			e.Field(i).SetUint(v)
		}

	}

	return &io, nil
}

func ReadProcessStatm(path string) (*ProcessStatm, error) {

	b, err := ioutil.ReadFile(path)

	if err != nil {
		return nil, err
	}

	s := string(b)
	f := strings.Fields(s)

	statm := ProcessStatm{}

	var n uint64

	for i := 0; i < len(f); i++ {

		if n, err = strconv.ParseUint(f[i], 10, 64); err != nil {
			return nil, err
		}

		switch i {
		case 0:
			statm.Size = n
		case 1:
			statm.Resident = n
		case 2:
			statm.Share = n
		case 3:
			statm.Text = n
		case 4:
			statm.Lib = n
		case 5:
			statm.Data = n
		case 6:
			statm.Dirty = n
		}

	}

	return &statm, nil
}

func ReadProcessStat(path string) (*ProcessStat, error) {

	b, err := ioutil.ReadFile(path)

	if err != nil {
		return nil, err
	}

	s := string(b)
	f := strings.Fields(s)

	stat := ProcessStat{}

	for i := 0; i < len(f); i++ {
		switch i {
		case 0:
			if stat.Pid, err = strconv.ParseUint(f[i], 10, 64); err != nil {
				return nil, err
			}
		case 1:
			stat.Comm = f[i]
		case 2:
			stat.State = f[i]
		case 3:
			if stat.Ppid, err = strconv.ParseInt(f[i], 10, 64); err != nil {
				return nil, err
			}
		case 4:
			if stat.Pgrp, err = strconv.ParseInt(f[i], 10, 64); err != nil {
				return nil, err
			}
		case 5:
			if stat.Session, err = strconv.ParseInt(f[i], 10, 64); err != nil {
				return nil, err
			}
		case 6:
			if stat.TtyNr, err = strconv.ParseInt(f[i], 10, 64); err != nil {
				return nil, err
			}
		case 7:
			if stat.Tpgid, err = strconv.ParseInt(f[i], 10, 64); err != nil {
				return nil, err
			}
		case 8:
			if stat.Flags, err = strconv.ParseUint(f[i], 10, 64); err != nil {
				return nil, err
			}
		case 9:
			if stat.Minflt, err = strconv.ParseUint(f[i], 10, 64); err != nil {
				return nil, err
			}
		case 10:
			if stat.Cminflt, err = strconv.ParseUint(f[i], 10, 64); err != nil {
				return nil, err
			}
		case 11:
			if stat.Majflt, err = strconv.ParseUint(f[i], 10, 64); err != nil {
				return nil, err
			}
		case 12:
			if stat.Cmajflt, err = strconv.ParseUint(f[i], 10, 64); err != nil {
				return nil, err
			}
		case 13:
			if stat.Utime, err = strconv.ParseUint(f[i], 10, 64); err != nil {
				return nil, err
			}
		case 14:
			if stat.Stime, err = strconv.ParseUint(f[i], 10, 64); err != nil {
				return nil, err
			}
		case 15:
			if stat.Cutime, err = strconv.ParseInt(f[i], 10, 64); err != nil {
				return nil, err
			}
		case 16:
			if stat.Cstime, err = strconv.ParseInt(f[i], 10, 64); err != nil {
				return nil, err
			}
		case 17:
			if stat.Priority, err = strconv.ParseInt(f[i], 10, 64); err != nil {
				return nil, err
			}
		case 18:
			if stat.Nice, err = strconv.ParseInt(f[i], 10, 64); err != nil {
				return nil, err
			}
		case 19:
			if stat.NumThreads, err = strconv.ParseInt(f[i], 10, 64); err != nil {
				return nil, err
			}
		case 20:
			if stat.Itrealvalue, err = strconv.ParseInt(f[i], 10, 64); err != nil {
				return nil, err
			}
		case 21:
			if stat.Starttime, err = strconv.ParseUint(f[i], 10, 64); err != nil {
				return nil, err
			}
		case 22:
			if stat.Vsize, err = strconv.ParseUint(f[i], 10, 64); err != nil {
				return nil, err
			}
		case 23:
			if stat.Rss, err = strconv.ParseInt(f[i], 10, 64); err != nil {
				return nil, err
			}
		case 24:
			if stat.Rsslim, err = strconv.ParseUint(f[i], 10, 64); err != nil {
				return nil, err
			}
		case 25:
			if stat.Startcode, err = strconv.ParseUint(f[i], 10, 64); err != nil {
				return nil, err
			}
		case 26:
			if stat.Endcode, err = strconv.ParseUint(f[i], 10, 64); err != nil {
				return nil, err
			}
		case 27:
			if stat.Startstack, err = strconv.ParseUint(f[i], 10, 64); err != nil {
				return nil, err
			}
		case 28:
			if stat.Kstkesp, err = strconv.ParseUint(f[i], 10, 64); err != nil {
				return nil, err
			}
		case 29:
			if stat.Kstkeip, err = strconv.ParseUint(f[i], 10, 64); err != nil {
				return nil, err
			}
		case 30:
			if stat.Signal, err = strconv.ParseUint(f[i], 10, 64); err != nil {
				return nil, err
			}
		case 31:
			if stat.Blocked, err = strconv.ParseUint(f[i], 10, 64); err != nil {
				return nil, err
			}
		case 32:
			if stat.Sigignore, err = strconv.ParseUint(f[i], 10, 64); err != nil {
				return nil, err
			}
		case 33:
			if stat.Sigcatch, err = strconv.ParseUint(f[i], 10, 64); err != nil {
				return nil, err
			}
		case 34:
			if stat.Wchan, err = strconv.ParseUint(f[i], 10, 64); err != nil {
				return nil, err
			}
		case 35:
			if stat.Nswap, err = strconv.ParseUint(f[i], 10, 64); err != nil {
				return nil, err
			}
		case 36:
			if stat.Cnswap, err = strconv.ParseUint(f[i], 10, 64); err != nil {
				return nil, err
			}
		case 37:
			if stat.ExitSignal, err = strconv.ParseInt(f[i], 10, 64); err != nil {
				return nil, err
			}
		case 38:
			if stat.Processor, err = strconv.ParseInt(f[i], 10, 64); err != nil {
				return nil, err
			}
		case 39:
			if stat.RtPriority, err = strconv.ParseUint(f[i], 10, 64); err != nil {
				return nil, err
			}
		case 40:
			if stat.Policy, err = strconv.ParseUint(f[i], 10, 64); err != nil {
				return nil, err
			}
		case 41:
			if stat.DelayacctBlkioTicks, err = strconv.ParseUint(f[i], 10, 64); err != nil {
				return nil, err
			}
		case 42:
			if stat.GuestTime, err = strconv.ParseUint(f[i], 10, 64); err != nil {
				return nil, err
			}
		case 43:
			if stat.CguestTime, err = strconv.ParseInt(f[i], 10, 64); err != nil {
				return nil, err
			}
		case 44:
			if stat.StartData, err = strconv.ParseUint(f[i], 10, 64); err != nil {
				return nil, err
			}
		case 45:
			if stat.EndData, err = strconv.ParseUint(f[i], 10, 64); err != nil {
				return nil, err
			}
		case 46:
			if stat.StartBrk, err = strconv.ParseUint(f[i], 10, 64); err != nil {
				return nil, err
			}
		case 47:
			if stat.ArgStart, err = strconv.ParseUint(f[i], 10, 64); err != nil {
				return nil, err
			}
		case 48:
			if stat.ArgEnd, err = strconv.ParseUint(f[i], 10, 64); err != nil {
				return nil, err
			}
		case 49:
			if stat.EnvStart, err = strconv.ParseUint(f[i], 10, 64); err != nil {
				return nil, err
			}
		case 50:
			if stat.EnvEnd, err = strconv.ParseUint(f[i], 10, 64); err != nil {
				return nil, err
			}
		case 51:
			if stat.ExitCode, err = strconv.ParseInt(f[i], 10, 64); err != nil {
				return nil, err
			}
		}
	}

	return &stat, nil
}

func ReadMaxPID(path string) (uint64, error) {

	b, err := ioutil.ReadFile(path)

	if err != nil {
		return 0, err
	}

	s := strings.TrimSpace(string(b))

	i, err := strconv.ParseUint(s, 10, 64)

	if err != nil {
		return 0, err
	}

	return i, nil

}

func ListPID(path string, max uint64) ([]uint64, error) {

	l := make([]uint64, 0, 5)

	for i := uint64(1); i <= max; i++ {

		p := filepath.Join(path, strconv.FormatUint(i, 10))

		s, err := os.Stat(p)

		if err != nil && !os.IsNotExist(err) {
			return nil, err
		}

		if err != nil || !s.IsDir() {
			continue
		}

		l = append(l, i)

	}

	return l, nil
}

// process info reader
// https://github.com/sharyanto/scripts/blob/master/explain-proc-stat
