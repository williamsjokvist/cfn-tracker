package cfn

import "testing"

func TestShouldEscalateToHeadful(t *testing.T) {
	tests := []struct {
		name     string
		result   loginResult
		headless bool
		want     bool
	}{
		{name: "needs human while headless", result: loginNeedsHuman, headless: true, want: true},
		{name: "already visible", result: loginNeedsHuman, headless: false, want: false},
		{name: "login succeeded", result: loginOK, headless: true, want: false},
		{name: "login failed", result: loginFailed, headless: true, want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := shouldEscalateToHeadful(tt.result, tt.headless); got != tt.want {
				t.Fatalf("shouldEscalateToHeadful() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestShouldReturnToHeadless(t *testing.T) {
	tests := []struct {
		name           string
		preferHeadless bool
		headless       bool
		want           bool
	}{
		{name: "return after visible authentication", preferHeadless: true, headless: false, want: true},
		{name: "already headless", preferHeadless: true, headless: true, want: false},
		{name: "visible debug mode", preferHeadless: false, headless: false, want: false},
		{name: "visible preferred but currently headless", preferHeadless: false, headless: true, want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := shouldReturnToHeadless(tt.preferHeadless, tt.headless); got != tt.want {
				t.Fatalf("shouldReturnToHeadless() = %v, want %v", got, tt.want)
			}
		})
	}
}
