package uuidv7

import (
	"testing"
	"time"
)

func TestNewV7(t *testing.T) {
	uuid := NewV7()
	if uuid == [16]byte{} {
		t.Error("NewV7() returned zero UUID")
	}
}

func TestNewV7Uniqueness(t *testing.T) {
	seen := make(map[UUIDv7]bool)
	for i := range 1000 {
		uuid := NewV7()
		if seen[uuid] {
			t.Errorf("Duplicate UUID generated at iteration %d", i)
		}
		seen[uuid] = true
	}
}

func TestNewV7Format(t *testing.T) {
	uuid := NewV7()

	// Check version bits (byte 6, bits 4-7 should be 0x70)
	version := (uuid[6] & 0xf0) >> 4
	if version != 0x7 {
		t.Errorf("Invalid version: got 0x%x, want 0x7", version)
	}

	// Check variant bits (byte 8, bits 6-7 should be 10)
	variant := (uuid[8] & 0xc0) >> 6
	if variant != 0x2 {
		t.Errorf("Invalid variant: got 0x%x, want 0x2", variant)
	}
}

func TestNewV7TimestampOrdering(t *testing.T) {
	uuid1 := NewV7()
	time.Sleep(1 * time.Millisecond)
	uuid2 := NewV7()

	// Extract timestamps (bytes 0-5, big-endian)
	ts1 := uint64(uuid1[0])<<40 | uint64(uuid1[1])<<32 | uint64(uuid1[2])<<24 | uint64(uuid1[3])<<16 | uint64(uuid1[4])<<8 | uint64(uuid1[5])
	ts2 := uint64(uuid2[0])<<40 | uint64(uuid2[1])<<32 | uint64(uuid2[2])<<24 | uint64(uuid2[3])<<16 | uint64(uuid2[4])<<8 | uint64(uuid2[5])

	if ts1 >= ts2 {
		t.Errorf("Timestamps not ordered: ts1=%d, ts2=%d", ts1, ts2)
	}
}

func TestGeneratorNew(t *testing.T) {
	gen := NewGenerator()
	uuid := gen.Next()
	if uuid == [16]byte{} {
		t.Error("Generator.New() returned zero UUID")
	}
}

func TestGeneratorUniqueness(t *testing.T) {
	gen := NewGenerator()
	seen := make(map[UUIDv7]bool)
	for i := range 1000 {
		uuid := gen.Next()
		if seen[uuid] {
			t.Errorf("Duplicate UUID generated at iteration %d", i)
		}
		seen[uuid] = true
	}
}

func TestMultipleGenerators(t *testing.T) {
	gen1 := NewGenerator()
	gen2 := NewGenerator()

	uuid1 := gen1.Next()
	uuid2 := gen2.Next()

	if uuid1 == uuid2 {
		t.Error("Different generators produced identical UUIDs")
	}
}

func TestGeneratorCounterSequencing(t *testing.T) {
	gen := NewGenerator()

	// Generate multiple UUIDs rapidly to test counter increment
	uuids := make([]UUIDv7, 100)
	for i := range 100 {
		uuids[i] = gen.Next()
	}

	// All should be unique
	seen := make(map[UUIDv7]bool)
	for i, uuid := range uuids {
		if seen[uuid] {
			t.Errorf("Duplicate UUID at index %d", i)
		}
		seen[uuid] = true
	}
}

func TestNewGeneratorWithBufferSize(t *testing.T) {
	gen := NewGeneratorWithBufferSize(16)
	uuid := gen.Next()
	if uuid == [16]byte{} {
		t.Error("NewGeneratorWithBufferSize() returned generator that produces zero UUID")
	}
}

func TestNewGeneratorWithBufferSizeMinimum(t *testing.T) {
	gen := NewGeneratorWithBufferSize(4)
	uuid := gen.Next()
	if uuid == [16]byte{} {
		t.Error("NewGeneratorWithBufferSize() with small size failed")
	}
}

func TestNewV7Concurrent(t *testing.T) {
	const numGoroutines = 100
	const numUUIDsPerGoroutine = 100

	results := make(chan UUIDv7, numGoroutines*numUUIDsPerGoroutine)

	for range numGoroutines {
		go func() {
			for range numUUIDsPerGoroutine {
				results <- NewV7()
			}
		}()
	}

	seen := make(map[UUIDv7]bool)
	for i := range numGoroutines * numUUIDsPerGoroutine {
		uuid := <-results
		if seen[uuid] {
			t.Errorf("Duplicate UUID generated concurrently at iteration %d", i)
		}
		seen[uuid] = true

		// Verify format
		version := (uuid[6] & 0xf0) >> 4
		if version != 0x7 {
			t.Errorf("Invalid version in concurrent UUID: got 0x%x, want 0x7", version)
		}
	}
}

func TestGeneratorConcurrent(t *testing.T) {
	const numGoroutines = 100
	const numUUIDsPerGoroutine = 100

	gen := NewGenerator()
	results := make(chan UUIDv7, numGoroutines*numUUIDsPerGoroutine)

	for range numGoroutines {
		go func() {
			for range numUUIDsPerGoroutine {
				results <- gen.Next()
			}
		}()
	}

	seen := make(map[UUIDv7]bool)
	for i := range numGoroutines * numUUIDsPerGoroutine {
		uuid := <-results
		if seen[uuid] {
			t.Errorf("Duplicate UUID generated concurrently at iteration %d", i)
		}
		seen[uuid] = true
	}
}
