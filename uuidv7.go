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
	defaultGenerator     *Generator
	defaultGeneratorOnce sync.Once
)

// Generator generates UUIDv7 values with counter-based sequencing and buffered random bytes
type Generator struct {
	mu                sync.Mutex
	lastTime          uint64
	clockSequence     uint16
	clockSequenceOnce sync.Once
	buffer            []byte
	cursor            int
}

// NewGenerator creates a new UUIDv7 generator with default buffer size (256 bytes = 32 UUIDs)
// 256 seems to be where diminishing returns start to kick in
func NewGenerator() *Generator {
	return NewGeneratorWithBufferSize(256)
}

// NewGeneratorWithBufferSize creates a new UUIDv7 generator with a custom buffer size
// Buffer size should be a multiple of 8 (bytes needed per UUID for random data)
func NewGeneratorWithBufferSize(bufferSize int) *Generator {
	if bufferSize < 8 {
		bufferSize = 8 // Minimum size
	}
	return &Generator{
		buffer: make([]byte, bufferSize),
		cursor: bufferSize, // Start at end to trigger initial fill
	}
}

// NewV7 generates a new UUIDv7 using the default generator
func NewV7() UUIDv7 {
	defaultGeneratorOnce.Do(func() {
		defaultGenerator = NewGenerator()
	})
	return defaultGenerator.NewV7()
}

// NewV7 generates a new UUIDv7
func (g *Generator) NewV7() UUIDv7 {
	g.mu.Lock()
	defer g.mu.Unlock()

	// Initialize counter randomly on first use
	g.clockSequenceOnce.Do(func() {
		buf := make([]byte, 2)
		_, _ = rand.Read(buf)
		g.clockSequence = binary.BigEndian.Uint16(buf)
	})

	// Get current time in milliseconds
	timeNow := uint64(time.Now().UnixMilli())

	// Increment counter if time hasn't changed
	if timeNow <= g.lastTime {
		g.clockSequence++
	}
	g.lastTime = timeNow

	// Get random bytes from buffer
	randBytes := g.nextRandBytes()

	// Build UUID
	var uuid UUIDv7

	// Bytes 0-5: 48-bit timestamp (big-endian)
	uuid[0] = byte(timeNow >> 40)
	uuid[1] = byte(timeNow >> 32)
	uuid[2] = byte(timeNow >> 24)
	uuid[3] = byte(timeNow >> 16)
	uuid[4] = byte(timeNow >> 8)
	uuid[5] = byte(timeNow)

	// Bytes 6-7: version (0x70) + sequence
	binary.BigEndian.PutUint16(uuid[6:8], g.clockSequence)
	uuid[6] = (uuid[6] & 0x0f) | 0x70 // Set version bits

	// Bytes 8-15: random bytes from buffer
	copy(uuid[8:16], randBytes)

	// Set variant bits (RFC 4122)
	uuid[8] = (uuid[8] & 0x3f) | 0x80

	return uuid
}

// nextRandBytes returns the next 8 random bytes from the buffer, refilling if needed
func (g *Generator) nextRandBytes() []byte {
	g.cursor += 8
	end := g.cursor + 8
	if end > len(g.buffer) {
		// Refill buffer
		_, _ = rand.Read(g.buffer)
		g.cursor = 0
		end = 8
	}
	return g.buffer[g.cursor:end]
}
