package uuid

import (
	"crypto/rand"
	"fmt"
)

// New creates a new UUID
func New() (string, error) {
	u := make([]byte, 16)
	_, err := rand.Read(u)
	if err != nil {
		return "", err
	}
	u[8] = u[8]&^0xc0 | 0x80
	u[6] = u[6]&^0xf0 | 0x40
	return fmt.Sprintf("%x-%x-%x-%x-%x", u[0:4], u[4:6], u[6:8], u[8:10], u[10:]), nil
}
