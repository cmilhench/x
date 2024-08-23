// Package env provides environment variable utilities.
package env_test

import (
	"testing"

	. "github.com/cmilhench/x/exp/env"
)

func TestGet(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name        string
		args        args
		shouldPanic bool
		want        string
	}{
		{
			name:        "should return value",
			args:        args{name: "TEST1"},
			shouldPanic: true,
		},
		{
			name:        "should return default value",
			args:        args{name: "TEST2"},
			shouldPanic: false,
			want:        "test2",
		},
	}
	t.Setenv("TEST2", "test2")
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.shouldPanic {
				if reason, isPanic := shouldPanic(t, func() { Get(tt.args.name) }); !isPanic {
					t.Errorf("Get() should panic, but got %v", reason)
				}
			} else {
				if got := Get(tt.args.name); got != tt.want {
					t.Errorf("Get() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

func TestGetDefault(t *testing.T) {
	type args struct {
		name  string
		value string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "should return value",
			args: args{name: "TEST1", value: "test1"},
			want: "test1",
		},
		{
			name: "should return default value",
			args: args{name: "TEST2", value: "default"},
			want: "test2",
		},
	}
	t.Setenv("TEST2", "test2")
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetDefault(tt.args.name, tt.args.value); got != tt.want {
				t.Errorf("GetDefault() = %v, want %v", got, tt.want)
			}
		})
	}
}

func shouldPanic(t *testing.T, f func()) (reason any, isPanic bool) {
	t.Helper()
	defer func() {
		if err := recover(); err != nil {
			reason = err
			isPanic = true
		}
	}()
	f()
	return nil, false
}
