package linux

import (
	"testing"
)

func TestMemInfo(t *testing.T) {
	info, err := ReadMemInfo("proc/meminfo")
	if err != nil {
		t.Fatal("meminfo read fail")
	}
	if info.Total == 0 {
		t.Fatal("total memory read fail")
	}
	if info.Free == 0 {
		t.Fatal("free memory read fail")
	}
	if info.Buffers == 0 {
		t.Fatal("buffers memory read fail")
	}
	t.Logf("%+v", info)
}
