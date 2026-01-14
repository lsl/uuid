// SPDX-License-Identifier: MPL-2.0
//
// Portions derived from:
//   - github.com/coolaj86/uuidv7 - Copyright 2024 AJ ONeal (MPL-2.0)
//   - github.com/gofrs/uuid - Copyright 2013–2018 Maxim Bublis (MIT)
//
// Modifications Copyright 2026 Louis Laugesen <louis@lgsn.dev>

package uuid

import (
	"crypto/rand"
	"encoding/binary"
	"sync"
	"time"
)

// UUIDv7 represents a UUIDv7 version 7 (16 bytes)
type UUIDv7 [16]byte

var (
	mu                sync.Mutex
	lastTime          uint64
	clockSequence     uint16
	clockSequenceOnce sync.Once
	buffer            []byte
	cursor            int
	bufferOnce        sync.Once
)

const defaultBufferSize = 256

// NewV7 generates a new UUIDv7
func NewV7() UUIDv7 {
	mu.Lock()
	defer mu.Unlock()

	// Initialize buffer on first use
	bufferOnce.Do(func() {
		buffer = make([]byte, defaultBufferSize)
		cursor = defaultBufferSize // Start at end to trigger initial fill
	})

	// Initialize counter randomly on first use
	clockSequenceOnce.Do(func() {
		buf := make([]byte, 2)
		_, _ = rand.Read(buf)
		clockSequence = binary.BigEndian.Uint16(buf)
	})

	// Get current time in milliseconds
	timeNow := uint64(time.Now().UnixMilli())

	// Increment counter if time hasn't changed
	if timeNow <= lastTime {
		clockSequence++
	}
	lastTime = timeNow

	// Get random bytes from buffer
	randBytes := nextRandBytes()

	// Build the UUIDv7
	var uuid UUIDv7

	// Bytes 0-5: 48-bit timestamp (big-endian)
	uuid[0] = byte(timeNow >> 40)
	uuid[1] = byte(timeNow >> 32)
	uuid[2] = byte(timeNow >> 24)
	uuid[3] = byte(timeNow >> 16)
	uuid[4] = byte(timeNow >> 8)
	uuid[5] = byte(timeNow)

	// Bytes 6-7: version (0x70) + sequence
	binary.BigEndian.PutUint16(uuid[6:8], clockSequence)
	uuid[6] = (uuid[6] & 0x0f) | 0x70 // Set version bits

	// Bytes 8-15: random bytes from buffer
	copy(uuid[8:16], randBytes)

	// Set variant bits (RFC 4122)
	uuid[8] = (uuid[8] & 0x3f) | 0x80

	return uuid
}

// nextRandBytes returns the next 8 random bytes from the buffer, refilling if needed
func nextRandBytes() []byte {
	cursor += 8
	end := cursor + 8
	if end > len(buffer) {
		// Refill buffer
		_, _ = rand.Read(buffer)
		cursor = 0
		end = 8
	}
	return buffer[cursor:end]
}
