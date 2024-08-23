package uuid

import (
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"time"
)

// New4 creates a new Version 4 UUID.
func New4() (string, error) {
	u := make([]byte, 16)
	_, err := rand.Read(u[:])
	if err != nil {
		return "", err
	}
	u[6] = (u[6] & 0x0f) | 0x40
	u[8] = (u[8] & 0x3f) | 0x80
	return fmt.Sprintf("%x-%x-%x-%x-%x", u[0:4], u[4:6], u[6:8], u[8:10], u[10:]), nil
}

// New7 creates a new Version 7 UUID.
// This naive implementation does not track previous sequence numbers, multiple
// IDs that are generated within the same millisecond, while unique, may not be
// ordered correctly
func New7() (string, error) {
	u := make([]byte, 16)
	_, err := rand.Read(u[:])
	if err != nil {
		return "", err
	}

	now := time.Now()
	tms := now.UnixMilli()
	seq := (now.UnixNano() - (tms * 1e6)) >> 8 // remaining ns within the ms

	binary.BigEndian.PutUint64(u[:8], uint64(tms)<<16)
	u[6] = (byte(seq>>8) & 0x0f) | 0x70
	u[7] = byte(seq)
	u[8] = (u[8] & 0x3f) | byte(1)
	return fmt.Sprintf("%x-%x-%x-%x-%x", u[0:4], u[4:6], u[6:8], u[8:10], u[10:]), nil
}
