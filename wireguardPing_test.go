package main

import (
	"net"
	"testing"
)

func TestWireGuardPing(t *testing.T) {
	tests := []struct {
		ip   net.IP
		port int
	}{
		{net.ParseIP("162.159.192.1"), 1701},
		{net.ParseIP("162.159.192.1"), 2400},
		{net.ParseIP("162.159.192.1"), 500},
		{net.ParseIP("162.159.192.1"), 4500},
	}

	for _, tt := range tests {
		t.Run(tt.ip.String(), func(t *testing.T) {
			got, duration := WireGuardPing(tt.ip, tt.port)
			if got != true {
				t.Errorf("WireGuardPing(%v, %v) = %v, want %v", tt.ip.String(), tt.port, got, duration)
			} else {
				t.Logf("WireGuardPing(%v, %v) = %v, want %v", tt.ip.String(), tt.port, got, duration)
			}
		})
	}
}
