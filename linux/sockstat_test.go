package linux

import (
	"reflect"
	"testing"
)

func TestSockStat(t *testing.T) {
	expected := SockStat{
		SocketsUsed:  231,
		TCPInUse:     27,
		TCPOrphan:    1,
		TCPTimeWait:  23,
		TCPAllocated: 31,
		TCPMemory:    3,
		UDPInUse:     19,
		UDPMemory:    17,
		UDPLITEInUse: 0,
		RAWInUse:     0,
		FRAGInUse:    0,
		FRAGMemory:   0,
	}
	sockStat, err := ReadSockStat("proc/sockstat")
	if err != nil {
		t.Fatal("sockstat read fail", err)
	}

	t.Logf("%+v", sockStat)

	if !reflect.DeepEqual(*sockStat, expected) {
		t.Error("not equal to expected")
	}
}
