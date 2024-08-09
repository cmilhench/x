package identifiers

import (
	"time"
)

// Creator generated a new unique identifier with the following characteristics
// 41 bits = milliseconds from epoch (max:2199023255551 = ~69 years)
// 10 bits = shard (max:1024)
// 12 bits = auto-incrementing and wrapping index (max:4095) see %
func Creator(shard uint16) func() uint64 {
	e := int64(1577836800000) // time.Parse(time.RFC3339, "2020-01-01T00:00:00Z")
	l := time.Now().UnixMilli() - e
	i := shard % 1024
	s := 0
	time.Sleep(time.Millisecond)
	return func() uint64 {
		var hash uint64
		n := time.Now().UnixMilli() - e
		if n == l {
			s = (s + 1) % 4095 // don't overflow the last 12 bits
			if s == 0 {
				// we have overflowed so wait until the next millisecond
				for n <= l {
					n = time.Now().UnixMilli() - e
				}
			}
		} else {
			s = 0
		}
		l = n
		// set the first 41 bits of the uint64 by shifting N left by 22
		// 0111111111111111111111111111111111111111110000000000000000000000
		hash |= uint64(n) << 22
		// set the next 10 bits (%1024) of the uint64 by shifting I left by 12
		// 0000000000000000000000000000000000000000001111111111000000000000
		hash |= uint64(i) << 12
		// set the last 12 bits (%4095) of the uint64 by shifting S left by 0
		// 0000000000000000000000000000000000000000000000000000111111111111
		hash |= uint64(s)
		return hash
	}
}

// Parse an identifier into it's components
func Parse(hash uint64) (time.Time, uint64, uint64, error) {
	e := int64(1577836800000) // time.Parse(time.RFC3339, "2020-01-01T00:00:00Z")
	n := (hash << (1)) >> 23
	i := (hash << (42)) >> 54
	s := (hash << (52)) >> 52

	// add the milliseconds to the epoch time to get the actual time
	now := time.UnixMilli(e).Add(time.Millisecond * time.Duration(int64(n))).UTC()

	return now, i, s, nil
}
