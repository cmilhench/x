package context

import (
	"context"
	"reflect"
	"testing"
)

func TestWithUser(t *testing.T) {
	type User struct {
		Name string
	}
	type args struct {
		parent context.Context
		user   *User
	}
	tests := []struct {
		name string
		args args
		want *User
	}{
		{
			name: "test",
			args: args{
				parent: context.Background(),
				user:   &User{Name: "test"},
			},
			want: &User{Name: "test"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := WithValueOf(tt.args.parent, tt.args.user)
			if got, ok := ValueOf[User](ctx); ok && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithUser() = %v, want %v", got, tt.want)
			}
		})
	}
}
