package linux

import (
	"errors"
	"io/ioutil"
	"strconv"
	"strings"
)

// UnixSocket represents a single Unix domain socket entry from /proc/net/unix.
type UnixSocket struct {
	Inode    uint64 `json:"inode"`
	RefCnt   uint64 `json:"ref_cnt"`
	Protocol uint64 `json:"protocol"`
	Flags    uint64 `json:"flags"`
	Type     uint64 `json:"type"`
	State    uint64 `json:"state"`
	Path     string `json:"path"`
}

// UnixSockets represents a list of Unix domain sockets parsed from /proc/net/unix.
type UnixSockets struct {
	Sockets []UnixSocket `json:"sockets"`
}

// NetUnixDomainSockets preserves backwards compatibility with existing callers.
type NetUnixDomainSockets struct {
	Sockets []NetUnixDomainSocket `json:"sockets"`
}

// NetUnixDomainSocket preserves backwards compatibility with existing callers.
type NetUnixDomainSocket struct {
	Protocol uint64 `json:"protocol"`
	RefCount uint64 `json:"ref_count"`
	Flags    uint64 `json:"flags"`
	Type     uint64 `json:"type"`
	State    uint64 `json:"state"`
	Inode    uint64 `json:"inode"`
	Path     string `json:"path"`
}

// ReadUnixSockets parses the /proc/net/unix socket table.
func ReadUnixSockets(path string) (*UnixSockets, error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(b), "\n")
	unixSockets := &UnixSockets{}

	for _, rawLine := range lines {
		line := strings.TrimSpace(rawLine)
		if line == "" {
			continue
		}
		f := strings.Fields(line)
		if len(f) < 7 || strings.HasPrefix(f[0], "Num") {
			continue
		}

		socket := UnixSocket{}
		if socket.RefCnt, err = strconv.ParseUint(f[1], 16, 64); err != nil {
			return nil, errors.New("Cannot parse unix socket [invalid RefCnt]: " + f[1])
		}

		if socket.Protocol, err = strconv.ParseUint(f[2], 16, 64); err != nil {
			return nil, errors.New("Cannot parse unix socket [invalid Protocol]: " + f[2])
		}

		if socket.Flags, err = strconv.ParseUint(f[3], 16, 64); err != nil {
			return nil, errors.New("Cannot parse unix socket [invalid Flags]: " + f[3])
		}

		if socket.Type, err = strconv.ParseUint(f[4], 16, 64); err != nil {
			return nil, errors.New("Cannot parse unix socket [invalid Type]: " + f[4])
		}

		if socket.State, err = strconv.ParseUint(f[5], 16, 64); err != nil {
			return nil, errors.New("Cannot parse unix socket [invalid State]: " + f[5])
		}

		if socket.Inode, err = strconv.ParseUint(f[6], 10, 64); err != nil {
			return nil, errors.New("Cannot parse unix socket [invalid Inode]: " + f[6])
		}

		if len(f) >= 8 {
			socket.Path = strings.Join(f[7:], " ")
		}

		unixSockets.Sockets = append(unixSockets.Sockets, socket)
	}

	return unixSockets, nil
}

// ReadNetUnixDomainSockets parses /proc/net/unix for backwards compatibility.
func ReadNetUnixDomainSockets(fpath string) (*NetUnixDomainSockets, error) {
	b, err := ioutil.ReadFile(fpath)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(b), "\n")
	unixDomainSockets := &NetUnixDomainSockets{}

	for i := 1; i < len(lines); i++ {
		line := lines[i]
		f := strings.Fields(line)

		if len(f) < 7 {
			continue
		}

		socket := NetUnixDomainSocket{}
		if socket.RefCount, err = strconv.ParseUint(f[1], 16, 64); err != nil {
			return nil, errors.New("Cannot parse unix domain socket [invalid RefCount]: " + f[1])
		}

		if socket.Protocol, err = strconv.ParseUint(f[2], 10, 64); err != nil {
			return nil, errors.New("Cannot parse unix domain socket [invalid Protocol]: " + f[2])
		}

		if socket.Flags, err = strconv.ParseUint(f[3], 10, 64); err != nil {
			return nil, errors.New("Cannot parse unix domain socket [invalid Flags]: " + f[3])
		}

		if socket.Type, err = strconv.ParseUint(f[4], 10, 64); err != nil {
			return nil, errors.New("Cannot parse unix domain socket [invalid Type]: " + f[4])
		}

		if socket.State, err = strconv.ParseUint(f[5], 10, 64); err != nil {
			return nil, errors.New("Cannot parse unix domain socket [invalid State]: " + f[5])
		}

		if socket.Inode, err = strconv.ParseUint(f[6], 10, 64); err != nil {
			return nil, errors.New("Cannot parse unix domain socket [invalid Inode]: " + f[6])
		}

		if len(f) >= 8 {
			socket.Path = f[7]
		}
		unixDomainSockets.Sockets = append(unixDomainSockets.Sockets, socket)
	}
	return unixDomainSockets, nil
}
