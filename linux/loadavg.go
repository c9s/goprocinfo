package linux

import (
	"io/ioutil"
	"strconv"
	"strings"
)

type LoadAvg struct {
	Last1Min  float64
	Last5Min  float64
	Last15Min float64
}

func ReadLoadAvg(path string) (*LoadAvg, error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	content := string(b)
	fields := strings.Fields(content)
	loadavg := LoadAvg{}
	loadavg.Last1Min, _ = strconv.ParseFloat(fields[0], 64)
	loadavg.Last5Min, _ = strconv.ParseFloat(fields[1], 64)
	loadavg.Last15Min, _ = strconv.ParseFloat(fields[2], 64)
	return &loadavg, nil
}
