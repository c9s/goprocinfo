package goproc

import "testing"

func TestLoadAvg(t *testing.T) {
	loadavg, err := ReadLoadAvg("proc/loadavg")
	if err != nil {
		t.Fatal("read loadavg fail")
	}
	t.Logf("%+v", loadavg)
}
