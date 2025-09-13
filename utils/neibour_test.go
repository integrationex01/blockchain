package utils_test

import (
	"blockchain/utils"
	"reflect"
	"testing"
)

func TestFindNeibours(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		myHost    string
		myPort    uint16
		startPort uint16
		endPort   uint16
		startIP   uint8
		endIp     uint8
		want      []string
	}{
		// TODO: Add test cases.
		{
			name:      "Test1",
			myHost:    "127.0.0.1",
			myPort:    5000,
			startPort: 5000,
			endPort:   5005,
			startIP:   0,
			endIp:     0,
			want:      []string{"127.0.0.1:5001"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := utils.FindNeibours(tt.myHost, tt.myPort, tt.startPort, tt.endPort, tt.startIP, tt.endIp)
			// TODO: update the condition below to compare got with tt.want.
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FindNeibours() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetHost(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		want string
	}{
		// TODO: Add test cases.
		{
			name: "Test1",
			want: "192.168.1.2",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := utils.GetHost()
			// TODO: update the condition below to compare got with tt.want.
			if got != tt.want {
				t.Errorf("GetHost() = %v, want %v", got, tt.want)
			}
		})
	}
}
