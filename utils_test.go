package fst

import "testing"

func TestLongestPrefix(t *testing.T) {
	type args struct {
		a []byte
		b []byte
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{name: "1", args: args{a: []byte(""), b: []byte("123")}, want: 0},
		{name: "2", args: args{a: []byte("123"), b: []byte("123")}, want: 3},
		{name: "3", args: args{a: []byte("1234"), b: []byte("123")}, want: 3},
		{name: "4", args: args{a: []byte("231"), b: []byte("123")}, want: 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := LongestPrefix(tt.args.a, tt.args.b); got != tt.want {
				t.Errorf("LongestPrefix() = %v, want %v", got, tt.want)
			}
		})
	}
}
