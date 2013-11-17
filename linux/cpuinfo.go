package linux

import (
	"fmt"
	"io/ioutil"
	"regexp"
	"strconv"
	"strings"
)

type CPUInfo struct {
	Processors []Processor
}

type Processor struct {
	Id        int64
	VendorId  string
	Model     int64
	ModelName string
	Flags     []string
	Cores     int64
	MHz       float64
}

var cpuinfoRegExp = regexp.MustCompile("([^:]*?)\\s*:\\s*(.*)$")

func ReadCPUInfo(path string) (*CPUInfo, error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	content := string(b)
	lines := strings.Split(content, "\n")

	var cpuinfo = CPUInfo{}
	var processor = &Processor{}

	for _, line := range lines {
		var key string
		var value string

		if len(line) == 0 {
			// end of processor
			cpuinfo.Processors = append(cpuinfo.Processors, *processor)
			processor = &Processor{}
			continue
		}

		submatches := cpuinfoRegExp.FindStringSubmatch(line)
		key = submatches[1]
		value = submatches[2]

		fmt.Printf("'%s'\n", key)

		switch key {
		case "processor":
			processor.Id, _ = strconv.ParseInt(value, 10, 32)
		case "vendor_id":
			processor.VendorId = value
		case "model":
			processor.Model, _ = strconv.ParseInt(value, 10, 32)
		case "model name":
			processor.ModelName = value
		case "flags":
			processor.Flags = strings.Fields(value)
		case "cpu cores":
			processor.Cores, _ = strconv.ParseInt(value, 10, 32)
		case "cpu MHz":
			processor.MHz, _ = strconv.ParseFloat(value, 64)
		}
		/*
			processor	: 0
			vendor_id	: GenuineIntel
			cpu family	: 6
			model		: 26
			model name	: Intel(R) Xeon(R) CPU           L5520  @ 2.27GHz
		*/
	}
	return &cpuinfo, nil
}
