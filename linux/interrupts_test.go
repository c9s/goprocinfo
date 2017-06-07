package linux

import "testing"

func TestInterrupts(t *testing.T) {
	interrupts, err := ReadInterrupts("proc/interrupts")
	if err != nil {
		t.Fatal("interrupts read fail")
	}
	_ = interrupts
	t.Logf("%+v", interrupts)
}
