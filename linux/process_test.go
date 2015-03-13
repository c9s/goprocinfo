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
		t.Fatal("process statm read fail", err)
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

func TestReadProcessStat(t *testing.T) {

	stat, err := ReadProcessStat("proc/3323/stat")

	if err != nil {
		t.Fatal("process stat read fail", err)
	}

	expected := &ProcessStat{
		Pid:                 3323,
		Comm:                "(proftpd)",
		State:               "S",
		Ppid:                1,
		Pgrp:                3323,
		Session:             3323,
		TtyNr:               0,
		Tpgid:               -1,
		Flags:               4202816,
		Minflt:              1311,
		Cminflt:             57367,
		Majflt:              0,
		Cmajflt:             1,
		Utime:               23,
		Stime:               58,
		Cutime:              24,
		Cstime:              49,
		Priority:            20,
		Nice:                0,
		NumThreads:          1,
		Itrealvalue:         0,
		Starttime:           2789,
		Vsize:               16601088,
		Rss:                 522,
		Rsslim:              4294967295,
		Startcode:           134512640,
		Endcode:             135222176,
		Startstack:          3217552592,
		Kstkesp:             3217551836,
		Kstkeip:             4118799382,
		Signal:              0,
		Blocked:             0,
		Sigignore:           272633856,
		Sigcatch:            8514799,
		Wchan:               0,
		Nswap:               0,
		Cnswap:              0,
		ExitSignal:          17,
		Processor:           7,
		RtPriority:          0,
		Policy:              0,
		DelayacctBlkioTicks: 1,
		GuestTime:           0,
		CguestTime:          0,
	}

	if !reflect.DeepEqual(stat, expected) {
		t.Error("not equal to expected", expected)
	}

	t.Logf("%+v", stat)

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
		Stat: ProcessStat{
			Pid:                 3323,
			Comm:                "(proftpd)",
			State:               "S",
			Ppid:                1,
			Pgrp:                3323,
			Session:             3323,
			TtyNr:               0,
			Tpgid:               -1,
			Flags:               4202816,
			Minflt:              1311,
			Cminflt:             57367,
			Majflt:              0,
			Cmajflt:             1,
			Utime:               23,
			Stime:               58,
			Cutime:              24,
			Cstime:              49,
			Priority:            20,
			Nice:                0,
			NumThreads:          1,
			Itrealvalue:         0,
			Starttime:           2789,
			Vsize:               16601088,
			Rss:                 522,
			Rsslim:              4294967295,
			Startcode:           134512640,
			Endcode:             135222176,
			Startstack:          3217552592,
			Kstkesp:             3217551836,
			Kstkeip:             4118799382,
			Signal:              0,
			Blocked:             0,
			Sigignore:           272633856,
			Sigcatch:            8514799,
			Wchan:               0,
			Nswap:               0,
			Cnswap:              0,
			ExitSignal:          17,
			Processor:           7,
			RtPriority:          0,
			Policy:              0,
			DelayacctBlkioTicks: 1,
			GuestTime:           0,
			CguestTime:          0,
		},
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

	var expected = []uint64{3323}

	if !reflect.DeepEqual(list, expected) {
		t.Error("not equal to expected", expected)
	}

	t.Logf("%+v", list)

}
