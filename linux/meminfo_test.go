package linux

import (
	"testing"
)

func TestMemInfo(t *testing.T) {
	info, err := ReadMemInfo("proc/meminfo")
	if err != nil {
		t.Fatal("meminfo read fail")
	}
	t.Logf("%+v", info)
}
