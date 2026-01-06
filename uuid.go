package uuid

import "github.com/lsl/uuid/uuidv7"

// UUIDv7 represents a UUID version 7 (16 bytes)
type UUIDv7 = uuidv7.UUIDv7

// NewV7 generates a new UUID version 7 using the default generator
func NewV7() UUIDv7 {
	return uuidv7.NewV7()
}
