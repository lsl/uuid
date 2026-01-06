package uuid

import "testing"

func TestNewV7(t *testing.T) {
	uuid := NewV7()
	if uuid == [16]byte{} {
		t.Error("NewV7() returned zero UUID")
	}
}
