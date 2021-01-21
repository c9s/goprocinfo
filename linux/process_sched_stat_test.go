package linux

import (
	"reflect"
	"testing"
)

func TestReadProcessSchedStat(t *testing.T) {
	schedStat, err := ReadProcessSchedStat("proc/3323/schedstat")
	if err != nil {
		t.Fatal("process sched-stat read fail", err)
	}

	expected := &ProcessSchedStat{
		RunTime:      10148346876,
		RunqueueTime: 2977307715,
		RunPeriods:   39798,
	}

	if !reflect.DeepEqual(schedStat, expected) {
		t.Errorf("not equal to expected %+v", expected)
	}

	t.Logf("%+v", schedStat)
}
