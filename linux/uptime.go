package linux

import (
	"io/ioutil"
	"strconv"
	"strings"
)

type Uptime struct {
	Total float64
	Idle  float64
}

func ReadUptime(path string) (*Uptime, err) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	fields := strings.Fields(string(b))
	uptime := Uptime{}
	if uptime.Total, err = strconv.ParseFloat(fields[0], 10, 64); err != nil {
		return nil, err
	}
	if uptime.Idle, err = strconv.ParseFloat(fields[1], 10, 64); err != nil {
		return nil, err
	}
	return &uptime, nil
}
