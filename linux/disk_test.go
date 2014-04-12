package linux

import "testing"

func TestDisk(t *testing.T) {
	disk, err := ReadDisk("/")
	if err != nil {
		t.Fatal("disk read fail")
	}
	if disk.Free <= 0 {
		t.Log("no good")
	}
	t.Logf("%+v", disk)
}
