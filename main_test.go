package main

import (
	"testing"
)

func Test_ellipsis(t *testing.T) {
	type args struct {
		text string
		max  int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{name: "Test max=1", args: struct {
			text string
			max  int
		}{text: "This is a very looooooooong text and we try to cut it", max: 1}, want: "Th ... it"},
		{name: "Test max=9", args: struct {
			text string
			max  int
		}{text: "This is a very looooooooong text and we try to cut it", max: 9}, want: "Th ... it"},
		{name: "Test max=11", args: struct {
			text string
			max  int
		}{text: "This is a very looooooooong text and we try to cut it", max: 11}, want: "Thi ...  it"},
		{name: "Test max=12", args: struct {
			text string
			max  int
		}{text: "This is a very looooooooong text and we try to cut it", max: 12}, want: "Thi ...  it"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ellipsis(tt.args.text, tt.args.max); got != tt.want {
				t.Errorf("ellipsis() = %v, want %v", got, tt.want)
			}
		})
	}
}
