package util

import "testing"

func TestPrettyPrint(t *testing.T) {
	type args struct {
		arr [][]string
	}
	type testObj struct {
		name string
		args args
	}
	tests := make([]testObj, 1)
	data := testObj{
		name: "normal",
		args: args{
			arr: [][]string{
				{"te", "test", "sdf"},
				{"te1123", "test123123", "123123"},
			},
		},
	}
	tests = append(tests, data)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			PrettyPrint(tt.args.arr)
		})
	}
}
