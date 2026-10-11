package linux

import (
	"io/ioutil"
	"path/filepath"
	"reflect"
	"testing"
)

func TestReadProcessStatus(t *testing.T) {

	status, err := ReadProcessStatus("proc/3323/status")

	if err != nil {
		t.Fatal("process io read fail", err)
	}

	expected := &ProcessStatus{
		Name:                     "proftpd",
		State:                    "S (sleeping)",
		Tgid:                     3323,
		Pid:                      3323,
		PPid:                     1,
		TracerPid:                0,
		RealUid:                  0,
		EffectiveUid:             111,
		SavedSetUid:              0,
		FilesystemUid:            111,
		RealGid:                  65534,
		EffectiveGid:             65534,
		SavedSetGid:              65534,
		FilesystemGid:            65534,
		FDSize:                   32,
		Groups:                   []int64{2001, 65534},
		VmPeak:                   16216,
		VmSize:                   16212,
		VmLck:                    0,
		VmHWM:                    2092,
		VmRSS:                    2088,
		VmData:                   872,
		VmStk:                    272,
		VmExe:                    696,
		VmLib:                    9416,
		VmPTE:                    36,
		VmSwap:                   0,
		Threads:                  1,
		SigQLength:               0,
		SigQLimit:                12091,
		SigPnd:                   0,
		ShdPnd:                   0,
		SigBlk:                   0,
		SigIgn:                   272633856,
		SigCgt:                   6450965743,
		CapInh:                   0,
		CapPrm:                   18446744073709551615,
		CapEff:                   0,
		CapBnd:                   18446744073709551615,
		Seccomp:                  0,
		CpusAllowed:              []uint32{255},
		CpusAllowedList:          []uint32{0, 1, 2, 3, 4, 5, 6, 7},
		VoluntaryCtxtSwitches:    5899,
		NonvoluntaryCtxtSwitches: 26,
	}

	if !reflect.DeepEqual(status, expected) {
		t.Error("not equal to expected", expected)
	}

	t.Logf("%+v", status)

}

func TestParseCPUList(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []uint32
	}{
		{
			name:     "empty string",
			input:    "",
			expected: nil,
		},
		{
			name:     "whitespace only",
			input:    "   \t\n  ",
			expected: nil,
		},
		{
			name:     "single cpu zero",
			input:    "0",
			expected: []uint32{0},
		},
		{
			name:     "single cpu non-zero",
			input:    "42",
			expected: []uint32{42},
		},
		{
			name:     "multiple comma-separated cpus",
			input:    "0,2,4,6",
			expected: []uint32{0, 2, 4, 6},
		},
		{
			name:     "simple continuous range",
			input:    "0-3",
			expected: []uint32{0, 1, 2, 3},
		},
		{
			name:     "single value range",
			input:    "5-5",
			expected: []uint32{5},
		},
		{
			name:     "combination of singles and ranges",
			input:    "0,2-4,7,9-11",
			expected: []uint32{0, 2, 3, 4, 7, 9, 10, 11},
		},
		{
			name:     "range with stride colon",
			input:    "0-6:2",
			expected: []uint32{0, 2, 4, 6},
		},
		{
			name:     "range with stride slash",
			input:    "1-7/3",
			expected: []uint32{1, 4, 7},
		},
		{
			name:     "whitespace within tokens",
			input:    " 0 - 2 , 5 , 8 - 10:2 ",
			expected: []uint32{0, 1, 2, 5, 8, 10},
		},
		{
			name:     "large cpu ids",
			input:    "1024-1026",
			expected: []uint32{1024, 1025, 1026},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := ParseCPUList(tc.input)
			if err != nil {
				t.Fatalf("unexpected error for %q: %v", tc.input, err)
			}
			if !reflect.DeepEqual(actual, tc.expected) {
				t.Errorf("input %q: expected %v, got %v", tc.input, tc.expected, actual)
			}
		})
	}
}

func TestParseCPUList_Errors(t *testing.T) {
	invalidInputs := []struct {
		name  string
		input string
	}{
		{"range start greater than end", "5-2"},
		{"multiple hyphens in range", "1-2-3"},
		{"missing range start", "-3"},
		{"missing range end", "3-"},
		{"non-numeric single", "abc"},
		{"non-numeric range start", "abc-3"},
		{"non-numeric range end", "1-abc"},
		{"stride zero with colon", "0-4:0"},
		{"stride zero with slash", "0-4/0"},
		{"multiple colons in stride", "0-4:2:1"},
		{"non-numeric stride", "0-4:xyz"},
		{"integer overflow uint32", "4294967296"},
	}

	for _, tc := range invalidInputs {
		t.Run(tc.name, func(t *testing.T) {
			result, err := ParseCPUList(tc.input)
			if err == nil {
				t.Fatalf("expected error for input %q, got result %v", tc.input, result)
			}
		})
	}
}

func TestReadProcessStatus_CpusAllowedList(t *testing.T) {
	tempDir := t.TempDir()

	t.Run("status with valid Cpus_allowed_list", func(t *testing.T) {
		content := "Name:\ttestproc\nState:\tS (sleeping)\nPid:\t1234\nCpus_allowed:\t0f\nCpus_allowed_list:\t0-3,5\n"
		statusFile := filepath.Join(tempDir, "status_valid")
		if err := ioutil.WriteFile(statusFile, []byte(content), 0644); err != nil {
			t.Fatalf("failed to write test fixture: %v", err)
		}

		status, err := ReadProcessStatus(statusFile)
		if err != nil {
			t.Fatalf("ReadProcessStatus failed: %v", err)
		}

		expectedList := []uint32{0, 1, 2, 3, 5}
		if !reflect.DeepEqual(status.CpusAllowedList, expectedList) {
			t.Errorf("expected CpusAllowedList %v, got %v", expectedList, status.CpusAllowedList)
		}
	})

	t.Run("status with stride in Cpus_allowed_list", func(t *testing.T) {
		content := "Name:\ttestproc\nState:\tS (sleeping)\nPid:\t1234\nCpus_allowed_list:\t0-6:2\n"
		statusFile := filepath.Join(tempDir, "status_stride")
		if err := ioutil.WriteFile(statusFile, []byte(content), 0644); err != nil {
			t.Fatalf("failed to write test fixture: %v", err)
		}

		status, err := ReadProcessStatus(statusFile)
		if err != nil {
			t.Fatalf("ReadProcessStatus failed: %v", err)
		}

		expectedList := []uint32{0, 2, 4, 6}
		if !reflect.DeepEqual(status.CpusAllowedList, expectedList) {
			t.Errorf("expected CpusAllowedList %v, got %v", expectedList, status.CpusAllowedList)
		}
	})

	t.Run("status without Cpus_allowed_list field", func(t *testing.T) {
		content := "Name:\ttestproc\nState:\tS (sleeping)\nPid:\t1234\n"
		statusFile := filepath.Join(tempDir, "status_no_list")
		if err := ioutil.WriteFile(statusFile, []byte(content), 0644); err != nil {
			t.Fatalf("failed to write test fixture: %v", err)
		}

		status, err := ReadProcessStatus(statusFile)
		if err != nil {
			t.Fatalf("ReadProcessStatus failed: %v", err)
		}

		if status.CpusAllowedList != nil {
			t.Errorf("expected nil CpusAllowedList when field is omitted, got %v", status.CpusAllowedList)
		}
	})

	t.Run("status with malformed Cpus_allowed_list", func(t *testing.T) {
		content := "Name:\ttestproc\nState:\tS (sleeping)\nPid:\t1234\nCpus_allowed_list:\tinvalid-range\n"
		statusFile := filepath.Join(tempDir, "status_malformed")
		if err := ioutil.WriteFile(statusFile, []byte(content), 0644); err != nil {
			t.Fatalf("failed to write test fixture: %v", err)
		}

		_, err := ReadProcessStatus(statusFile)
		if err == nil {
			t.Fatal("expected error for malformed Cpus_allowed_list, got nil")
		}
	})
}
