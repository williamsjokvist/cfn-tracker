package browser

import (
	"testing"

	"github.com/go-rod/rod/lib/proto"
)

func TestShouldBlockRequest(t *testing.T) {
	tests := []struct {
		name     string
		typeName proto.NetworkResourceType
		hostname string
		want     bool
	}{
		{"image", proto.NetworkResourceTypeImage, "example.com", true},
		{"font", proto.NetworkResourceTypeFont, "example.com", true},
		{"stylesheet", proto.NetworkResourceTypeStylesheet, "example.com", true},
		{"steam stylesheet", proto.NetworkResourceTypeStylesheet, "steamcommunity.com", false},
		{"script", proto.NetworkResourceTypeScript, "example.com", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := shouldBlockRequest(tt.typeName, tt.hostname); got != tt.want {
				t.Fatalf("shouldBlockRequest() = %v, want %v", got, tt.want)
			}
		})
	}
}
