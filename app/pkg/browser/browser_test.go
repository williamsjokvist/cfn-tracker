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
		blocking bool
		want     bool
	}{
		{"image", proto.NetworkResourceTypeImage, "example.com", true, true},
		{"font", proto.NetworkResourceTypeFont, "example.com", true, true},
		{"stylesheet", proto.NetworkResourceTypeStylesheet, "example.com", true, true},
		{"steam stylesheet", proto.NetworkResourceTypeStylesheet, "steamcommunity.com", true, false},
		{"script", proto.NetworkResourceTypeScript, "example.com", true, false},
		{"image blocking disabled", proto.NetworkResourceTypeImage, "example.com", false, false},
		{"font blocking disabled", proto.NetworkResourceTypeFont, "example.com", false, false},
		{"stylesheet blocking disabled", proto.NetworkResourceTypeStylesheet, "example.com", false, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := shouldBlockRequest(tt.typeName, tt.hostname, tt.blocking); got != tt.want {
				t.Fatalf("shouldBlockRequest() = %v, want %v", got, tt.want)
			}
		})
	}
}
