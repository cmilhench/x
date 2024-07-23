package ptr

import "time"

func Int32(i int32) *int32 {
	return &i
}

func Int32Value(i *int32) int32 {
	if i != nil {
		return *i
	}
	return 0
}

func Int64(i int64) *int64 {
	return &i
}

func Int64Value(i *int64) int64 {
	if i != nil {
		return *i
	}
	return 0
}

func String(s string) *string {
	return &s
}

func StringValue(s *string) string {
	if s != nil {
		return *s
	}
	return ""
}

func Bool(b bool) *bool {
	return &b
}

func BoolValue(b *bool) bool {
	if b != nil {
		return *b
	}
	return false
}

func Float32(f float32) *float32 {
	return &f
}

func Float32Value(f *float32) float32 {
	if f != nil {
		return *f
	}
	return 0
}

func Float64(f float64) *float64 {
	return &f
}

func Float64Value(f *float64) float64 {
	if f != nil {
		return *f
	}
	return 0
}

func Time(t time.Time) *time.Time {
	return &t
}

func TimeValue(t *time.Time) time.Time {
	if t != nil {
		return *t
	}
	return time.Time{}
}
