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

type ProcessStat struct{}

func ReadProcess(pid uint64, path string) (*Process, error) {

	p := filepath.Join(path, strconv.FormatUint(pid, 10))

	if _, err := os.Stat(p); err != nil {
		return nil, err
	}

	process := Process{}

	statPath := filepath.Join(p, "stat")
	statusPath := filepath.Join(p, "status")

	var err error
	var io *ProcessIO
	var statm *ProcessStatm

	if io, err = ReadProcessIO(filepath.Join(p, "io")); err != nil {
		return nil, err
	}

	if statm, err = ReadProcessStatm(filepath.Join(p, "statm")); err != nil {
		return nil, err
	}

	_ = statPath
	_ = statusPath

	process.IO = *io
	process.Statm = *statm

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
		v, err := strconv.ParseUint(l[1], 10, 0)

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

	for i := 0; i < 7; i++ {

		if n, err = strconv.ParseUint(f[i], 10, 0); err != nil {
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

func ReadMaxPID(path string) (uint64, error) {

	b, err := ioutil.ReadFile(path)

	if err != nil {
		return 0, err
	}

	s := strings.TrimSpace(string(b))

	i, err := strconv.ParseUint(s, 10, 0)

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
