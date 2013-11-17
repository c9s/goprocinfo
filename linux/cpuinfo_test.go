package linux

import "testing"

func TestCPUInfo(t *testing.T) {
	cpuinfo, err := ReadCPUInfo("proc/cpuinfo")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v", cpuinfo)

	if len(cpuinfo.Processors) != 8 {
		t.Fatal("wrong processor number")
	}
}
