package model

import (
	"testing"
)

func TestNewPair(t *testing.T) {
	type args struct {
		key   string
		value string
	}
	tests := []struct {
		name    string
		args    args
		want    *Pair
		wantErr bool
	}{
		{"test1", args{"5", "**"}, nil, true},
		{"test2", args{"5**", "a"}, nil, true},
		{"test3", args{"*$@#$@5", "a"}, nil, true},
		{"test4", args{"5", "a**"}, nil, true},
		{"test5", args{"5", "a"}, &Pair{Key: "5", Value: "a"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewPair(tt.args.key, tt.args.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewPair() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != nil {
				if got.Key != tt.want.Key || got.Value != tt.want.Value {
					t.Errorf("NewPair() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}
