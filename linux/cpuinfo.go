package linux

import "io/ioutil"

type CPUInfo struct {
	Processors []Processor
}

type Processor struct {
}

func ReadCPUInfo(path string) (*CPUInfo, error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	cpuinfo := CPUInfo{}

	return &cpuinfo, nil
}
