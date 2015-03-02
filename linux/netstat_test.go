package linux

import "testing"
import "reflect"

func TestNetstat(t *testing.T) {

	var expected = Netstat{0, 0, 1764, 180, 0, 0, 0, 0, 0, 0, 28321, 0, 0, 0, 0, 243, 25089, 53, 837, 0, 0, 95994, 623148353, 640988091, 0, 92391, 81263, 594305, 590571, 35, 6501, 81, 113, 213, 1, 223, 318, 1056, 287, 218, 6619, 435, 1, 975, 264, 17298, 871, 5836, 3843, 2, 520, 0, 0, 833, 0, 3235, 44, 0, 571, 163, 0, 138, 0, 0, 0, 19, 1312, 677, 129, 0, 0, 27986, 27713, 40522, 837, 0, 38648, 0, 0, 0, 0, 0, 0, 0, 0, 2772402103, 5189844022, 0, 0, 0, 0}

	netstat, err := ReadNetstat("proc/net_netstat")
	if err != nil {
		t.Fatal("netstat read fail", err)
	}

	t.Logf("%+v", netstat)

	if !reflect.DeepEqual(*netstat, expected) {
		t.Error("not equal to expected")
	}
}
