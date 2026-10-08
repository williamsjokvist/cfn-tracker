package cfn

import "testing"

func TestShouldEscalateToManualLogin(t *testing.T) {
	tests := []struct {
		name   string
		result loginResult
		want   bool
	}{
		// Chrome under rod can't pass Cloudflare even headful, so always fall back
		// to manual login.
		{name: "needs human", result: loginNeedsHuman, want: true},
		{name: "login succeeded", result: loginOK, want: false},
		{name: "login failed", result: loginFailed, want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := shouldEscalateToManualLogin(tt.result); got != tt.want {
				t.Fatalf("shouldEscalateToManualLogin() = %v, want %v", got, tt.want)
			}
		})
	}
}
