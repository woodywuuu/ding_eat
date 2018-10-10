package main

import (
	"testing"
)

func Test_get_weekday(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"1", args{"Monday"}, "周一"},
		{"2", args{"Tuesday"}, "周二"},
		{"3", args{"Wednesday"}, "周三"},
		{"4", args{"Thursday"}, "周四"},
		{"5", args{"Friday"}, "周五"},
		{"6", args{"Saturday"}, "周六"},
		{"7", args{"Sunday"}, "周日"},
		{"8", args{""}, ""},
		{"9", args{"abc"}, ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := get_weekday(tt.args.name); got != tt.want {
				t.Errorf("get_weekday() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_get_message(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name  string
		args  args
		want  string
		want1 string
	}{
		{"1", args{"周一"}, "周一", "么么叽，订外卖啦~~"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := get_message(tt.args.name)
			if got != tt.want {
				t.Errorf("get_message() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("get_message() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

