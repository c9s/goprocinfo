package linux

import (
	"reflect"
	"testing"
)

func TestReadProcessIO(t *testing.T) {

	io, err := ReadProcessIO("proc/3323/io")

	if err != nil {
		t.Fatal("process io read fail", err)
	}

	expected := &ProcessIO{
		RChar:               3865585,
		WChar:               183294,
		Syscr:               6697,
		Syscw:               997,
		ReadBytes:           90112,
		WriteBytes:          45056,
		CancelledWriteBytes: 0,
	}

	if !reflect.DeepEqual(io, expected) {
		t.Error("not equal to expected", expected)
	}

	t.Logf("%+v", io)
}

func TestReadProcessStatm(t *testing.T) {

	statm, err := ReadProcessStatm("proc/3323/statm")

	if err != nil {
		t.Fatal("process io read fail", err)
	}

	expected := &ProcessStatm{
		Size:     4053,
		Resident: 522,
		Share:    174,
		Text:     174,
		Lib:      0,
		Data:     286,
		Dirty:    0,
	}

	if !reflect.DeepEqual(statm, expected) {
		t.Error("not equal to expected", expected)
	}

	t.Logf("%+v", statm)

}

func TestReadProcess(t *testing.T) {

	p, err := ReadProcess(3323, "proc")

	if err != nil {
		t.Fatal("process read fail", err)
	}

	expected := &Process{

		Status: ProcessStatus{},
		Statm: ProcessStatm{
			Size:     4053,
			Resident: 522,
			Share:    174,
			Text:     174,
			Lib:      0,
			Data:     286,
			Dirty:    0,
		},
		Stat: ProcessStat{},

		IO: ProcessIO{
			RChar:               3865585,
			WChar:               183294,
			Syscr:               6697,
			Syscw:               997,
			ReadBytes:           90112,
			WriteBytes:          45056,
			CancelledWriteBytes: 0,
		},
	}

	if !reflect.DeepEqual(p, expected) {
		t.Error("not equal to expected", expected)
	}

	t.Logf("%+v", p)
}

func TestMaxPID(t *testing.T) {

	max, err := ReadMaxPID("proc/sys_kernel_pid_max")

	if err != nil {
		t.Fatal("max pid read fail", err)
	}

	_ = max

	if max != 32768 {
		t.Error("unexpected value")
	}

	t.Logf("%+v", max)

}

func TestListPID(t *testing.T) {

	list, err := ListPID("proc", 32768)

	if err != nil {
		t.Fatal("list pid fail", err)
	}

	_ = list

	var expected = []uint64{3323}

	if !reflect.DeepEqual(list, expected) {
		t.Error("not equal to expected", expected)
	}

	t.Logf("%+v", list)

}
