package linux

import (
	"reflect"
	"testing"
)

func TestReadNetUnixDomainSockets(t *testing.T) {
	type args struct {
		fpath string
	}
	tests := []struct {
		name string
		args func(t *testing.T) args

		want1      *NetUnixDomainSockets
		wantErr    bool
		inspectErr func(err error, t *testing.T) //use for more precise error evaluation after test
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
