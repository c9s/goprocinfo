package linux

import (
	"io/ioutil"
	"os"
	"reflect"
	"testing"
)

func TestReadNetUnixDomainSockets(t *testing.T) {
	type args struct {
		fpath string
	}
	tests := []struct {
		name       string
		args       func(t *testing.T) args
		want1      *NetUnixDomainSockets
		wantErr    bool
		inspectErr func(err error, t *testing.T)
	}{
		{
			name: "Read example file without errors",
			args: func(t *testing.T) args {
				return args{
					fpath: "proc/net_unix",
				}
			},
			want1:      expectedData(),
			wantErr:    false,
			inspectErr: func(err error, t *testing.T) {},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tArgs := tt.args(t)

			got1, err := ReadNetUnixDomainSockets(tArgs.fpath)
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("ReadNetUnixDomainSockets got1 = %v, want1: %v", got1, tt.want1)
			}

			if (err != nil) != tt.wantErr {
				t.Fatalf("ReadNetUnixDomainSockets error = %v, wantErr: %t", err, tt.wantErr)
			}

			if tt.inspectErr != nil {
				tt.inspectErr(err, t)
			}
		})
	}
}

func expectedData() *NetUnixDomainSockets {
	sockets := []NetUnixDomainSocket{
		NetUnixDomainSocket{
			Protocol: 0,
			RefCount: 2,
			Flags:    10000,
			Type:     1,
			State:    1,
			Inode:    10748,
			Path:     "/var/run/vmware/guestServicePipe",
		},
		NetUnixDomainSocket{
			Protocol: 0,
			RefCount: 2,
			Flags:    10000,
			Type:     1,
			State:    1,
			Inode:    53314618,
			Path:     "/tmp/nvimnhGCOs/0",
		},
		NetUnixDomainSocket{
			Protocol: 0,
			RefCount: 2,
			Flags:    10000,
			Type:     1,
			State:    1,
			Inode:    15567,
			Path:     "/var/run/nslcd/socket",
		},
		NetUnixDomainSocket{
			Protocol: 0,
			RefCount: 2,
			Flags:    10000,
			Type:     1,
			State:    1,
			Inode:    15689,
			Path:     "/var/run/rpcbind.sock",
		},
		NetUnixDomainSocket{
			Protocol: 0,
			RefCount: 2,
			Flags:    10000,
			Type:     1,
			State:    1,
			Inode:    445862525,
			Path:     "/tmp/nvimJiXQoJ/0",
		},
		NetUnixDomainSocket{
			Protocol: 0,
			RefCount: 2,
			Flags:    10000,
			Type:     1,
			State:    1,
			Inode:    215283,
			Path:     "/tmp/tmux-9202/default",
		},
		NetUnixDomainSocket{
			Protocol: 0,
			RefCount: 2,
			Flags:    10000,
			Type:     1,
			State:    1,
			Inode:    16303,
			Path:     "/var/run/dbus/system_bus_socket",
		},
		NetUnixDomainSocket{
			Protocol: 0,
			RefCount: 2,
			Flags:    10000,
			Type:     1,
			State:    1,
			Inode:    16721,
			Path:     "/var/run/cgred.socket",
		},
		NetUnixDomainSocket{
			Protocol: 0,
			RefCount: 2,
			Flags:    10000,
			Type:     1,
			State:    1,
			Inode:    16804,
			Path:     "/var/run/nscd/socket",
		},
	}

	uds := &NetUnixDomainSockets{}
	uds.Sockets = sockets
	return uds
}

func TestReadUnixSockets_ProcNetUnix(t *testing.T) {
	got, err := ReadUnixSockets("proc/net_unix")
	if err != nil {
		t.Fatalf("ReadUnixSockets failed: %v", err)
	}

	if len(got.Sockets) != 9 {
		t.Fatalf("expected 9 sockets, got %d", len(got.Sockets))
	}

	first := got.Sockets[0]
	if first.Inode != 10748 {
		t.Errorf("expected Inode 10748, got %d", first.Inode)
	}
	if first.RefCnt != 2 {
		t.Errorf("expected RefCnt 2, got %d", first.RefCnt)
	}
	if first.Protocol != 0 {
		t.Errorf("expected Protocol 0, got %d", first.Protocol)
	}
	if first.Flags != 0x00010000 {
		t.Errorf("expected Flags 0x10000, got %x", first.Flags)
	}
	if first.Type != 1 {
		t.Errorf("expected Type 1, got %d", first.Type)
	}
	if first.State != 1 {
		t.Errorf("expected State 1, got %d", first.State)
	}
	if first.Path != "/var/run/vmware/guestServicePipe" {
		t.Errorf("expected Path /var/run/vmware/guestServicePipe, got %s", first.Path)
	}
}

func TestReadUnixSockets_MultipleEntries(t *testing.T) {
	sampleContent := `Num       RefCount Protocol Flags    Type St Inode Path
0000000000000001: 00000002 00000000 00010000 0001 0A 12001 /run/systemd/journal/stdout
0000000000000002: 00000003 00000000 00000000 0001 01 12002 /tmp/app.sock
0000000000000003: 00000002 00000000 00000000 0001 01 12003
0000000000000004: 00000002 00000000 00000000 0002 07 12004 @/run/systemd/notify
0000000000000005: 00000001 00000000 00000000 0002 01 12005
0000000000000006: 00000002 00000000 00010000 0005 0A 12006 /tmp/seqpacket.sock
`
	tmpFile, err := ioutil.TempFile("", "proc_net_unix_test_*")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.WriteString(sampleContent); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}
	tmpFile.Close()

	sockets, err := ReadUnixSockets(tmpFile.Name())
	if err != nil {
		t.Fatalf("ReadUnixSockets failed: %v", err)
	}

	expected := &UnixSockets{
		Sockets: []UnixSocket{
			{
				Inode:    12001,
				RefCnt:   2,
				Protocol: 0,
				Flags:    0x00010000,
				Type:     1, // SOCK_STREAM
				State:    10, // TCP_LISTEN (listening)
				Path:     "/run/systemd/journal/stdout",
			},
			{
				Inode:    12002,
				RefCnt:   3,
				Protocol: 0,
				Flags:    0,
				Type:     1, // SOCK_STREAM
				State:    1, // TCP_ESTABLISHED (connected)
				Path:     "/tmp/app.sock",
			},
			{
				Inode:    12003,
				RefCnt:   2,
				Protocol: 0,
				Flags:    0,
				Type:     1, // SOCK_STREAM
				State:    1, // connected, unbound/anonymous
				Path:     "",
			},
			{
				Inode:    12004,
				RefCnt:   2,
				Protocol: 0,
				Flags:    0,
				Type:     2, // SOCK_DGRAM
				State:    7, // unconnected
				Path:     "@/run/systemd/notify",
			},
			{
				Inode:    12005,
				RefCnt:   1,
				Protocol: 0,
				Flags:    0,
				Type:     2, // SOCK_DGRAM
				State:    1, // connected, unbound/anonymous
				Path:     "",
			},
			{
				Inode:    12006,
				RefCnt:   2,
				Protocol: 0,
				Flags:    0x00010000,
				Type:     5, // SOCK_SEQPACKET
				State:    10, // listening
				Path:     "/tmp/seqpacket.sock",
			},
		},
	}

	if !reflect.DeepEqual(sockets, expected) {
		t.Errorf("sockets do not match expected.\nGot:  %+v\nWant: %+v", sockets, expected)
	}
}

func TestReadUnixSockets_EdgeCases(t *testing.T) {
	testCases := []struct {
		name    string
		content string
		wantLen int
	}{
		{
			name:    "empty file",
			content: "",
			wantLen: 0,
		},
		{
			name:    "header only",
			content: "Num       RefCount Protocol Flags    Type St Inode Path\n",
			wantLen: 0,
		},
		{
			name: "blank lines and partial lines",
			content: `
Num       RefCount Protocol Flags    Type St Inode Path

short line

0000000000000001: 00000001 00000000 00000000 0001 01 9999 /tmp/test.sock

`,
			wantLen: 1,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tmpFile, err := ioutil.TempFile("", "proc_net_unix_edge_*")
			if err != nil {
				t.Fatalf("failed to create temp file: %v", err)
			}
			defer os.Remove(tmpFile.Name())

			if _, err := tmpFile.WriteString(tc.content); err != nil {
				t.Fatalf("failed to write temp file: %v", err)
			}
			tmpFile.Close()

			res, err := ReadUnixSockets(tmpFile.Name())
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if len(res.Sockets) != tc.wantLen {
				t.Errorf("expected %d sockets, got %d", tc.wantLen, len(res.Sockets))
			}
		})
	}
}

func TestReadUnixSockets_Errors(t *testing.T) {
	// Non-existent file
	_, err := ReadUnixSockets("/path/does/not/exist/unix")
	if err == nil {
		t.Error("expected error for non-existent file, got nil")
	}

	invalidCases := []struct {
		name    string
		content string
	}{
		{
			name:    "invalid RefCnt",
			content: "Num RefCount Protocol Flags Type St Inode Path\n01: invalid_ref 00000000 00000000 0001 01 100 /tmp/s\n",
		},
		{
			name:    "invalid Protocol",
			content: "Num RefCount Protocol Flags Type St Inode Path\n01: 00000001 invalid_proto 00000000 0001 01 100 /tmp/s\n",
		},
		{
			name:    "invalid Flags",
			content: "Num RefCount Protocol Flags Type St Inode Path\n01: 00000001 00000000 invalid_flags 0001 01 100 /tmp/s\n",
		},
		{
			name:    "invalid Type",
			content: "Num RefCount Protocol Flags Type St Inode Path\n01: 00000001 00000000 00000000 invalid_type 01 100 /tmp/s\n",
		},
		{
			name:    "invalid State",
			content: "Num RefCount Protocol Flags Type St Inode Path\n01: 00000001 00000000 00000000 0001 invalid_st 100 /tmp/s\n",
		},
		{
			name:    "invalid Inode",
			content: "Num RefCount Protocol Flags Type St Inode Path\n01: 00000001 00000000 00000000 0001 01 invalid_inode /tmp/s\n",
		},
	}

	for _, ic := range invalidCases {
		t.Run(ic.name, func(t *testing.T) {
			tmpFile, err := ioutil.TempFile("", "proc_net_unix_err_*")
			if err != nil {
				t.Fatalf("failed to create temp file: %v", err)
			}
			defer os.Remove(tmpFile.Name())

			if _, err := tmpFile.WriteString(ic.content); err != nil {
				t.Fatalf("failed to write temp file: %v", err)
			}
			tmpFile.Close()

			_, err = ReadUnixSockets(tmpFile.Name())
			if err == nil {
				t.Errorf("expected error for case %s, got nil", ic.name)
			}
		})
	}
}
