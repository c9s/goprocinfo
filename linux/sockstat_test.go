package linux

import "testing"
import "reflect"

func TestSockstat(t *testing.T) {
	var expected = Sockstat{231, 27, 1, 23, 31, 3, 19, 17, 0, 0, 0, 0}

	sockstat, err := ReadSockstat("proc/sockstat")
	if err != nil {
		t.Fatal("sockstat read fail", err)
	}

	t.Logf("%+v", sockstat)

	if !reflect.DeepEqual(*sockstat, expected) {
		t.Error("not equal to expected")
	}
}
