package linux

import (
	"io/ioutil"
	"strconv"
	"strings"
)

/*
type MemInfo struct {
	Total      uint64
	Free       uint64
	Buffers    uint64
	Cached     uint64
	SwapCached uint64
	Active     uint64
	Inactive   uint64
}
*/
type MemInfo map[string]uint64

func parseMemInfo(content string) MemInfo {
	lines := strings.Split(content, "\n")
	var info = MemInfo{}
	for _, line := range lines {
		fields := strings.SplitN(line, ":", 2)
		if len(fields) < 2 {
			continue
		}
		keyField := fields[0]
		valFields := strings.Fields(fields[1])
		val, _ := strconv.ParseUint(valFields[0], 10, 64)
		info[keyField] = val
	}
	return info
}

func ReadMemInfo(path string) (MemInfo, error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	info := parseMemInfo(string(b))
	return info, nil
}
