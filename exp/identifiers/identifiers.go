package identifiers

import (
	"strconv"
	"time"
)

// Creator generated a new unique identifier with the following characteristics
// 41 bits = milliseconds from epoch (max:2199023255551 = ~69 years)
// 13 bits = shard (max:8191)
// 10 bits = auto-incrementing and wrapping index (max:1023) see %
func Creator(shard uint16) func() string {
	e, _ := time.Parse(time.RFC3339, "2020-01-01T00:00:00Z")
	s := shard
	i := 0
	return func() string {
		var hash uint64
		i = (i + 1) % 1024 // don't overflow the last 10 bits
		d := time.Since(e).Milliseconds()
		// set the first 41 bits (0:41)
		hash |= uint64(d) << (64 - 0 - 41) // << 23
		// set the next 13 bits (41:13)
		hash |= uint64(s) << (64 - 41 - 13) // << 10
		// set the last 10 bits (54:10)
		hash |= uint64(i) << (64 - 54 - 10)
		return strconv.FormatUint(hash, 36)
	}
}

// Parse an identifier into it's components
func Parse(key string) (time.Time, uint64, uint64, error) {
	e, _ := time.Parse(time.RFC3339, "2020-01-01T00:00:00Z")
	hash, err := strconv.ParseUint(key, 36, 64)
	if err != nil {
		return time.Time{}, 0, 0, err
	}
	d := (hash << (0)) >> (64 - 41)
	s := (hash << (41)) >> (64 - 13)
	i := (hash << (54)) >> (64 - 10)
	now := e.Add(time.Millisecond * time.Duration(int64(d)))
	return now, s, i, nil
}
