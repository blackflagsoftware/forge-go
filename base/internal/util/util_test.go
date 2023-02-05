package util

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetTypeCount(t *testing.T) {
	type args struct {
		i interface{}
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			"map success",
			args{map[string]int{
				"one": 1,
				"two": 2,
			}},
			2,
		},
		{
			"string success",
			args{"my string"},
			1,
		},
		{
			"array success",
			args{[2]int{1, 2}},
			2,
		},
		{
			"slice success",
			args{make([]int, 2, 5)},
			2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetTypeCount(tt.args.i); got != tt.want {
				t.Errorf("GetTypeCount() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValidJson(t *testing.T) {
	tests := []struct {
		name      string
		jsonValue json.RawMessage
		want      bool
	}{
		{
			"successful",
			[]byte(`{"test": "again"}`),
			true,
		},
		{
			"failed",
			[]byte(`{"test": again"}`),
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ValidJson(tt.jsonValue)
			assert.Equal(t, tt.want, got, "ValidJson.%s => expected %t; got %t", tt.name, tt.want, got)
		})
	}
}
