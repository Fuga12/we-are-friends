package httpapi

import "testing"

func TestValidAdminCredentials(t *testing.T) {
	tests := []struct {
		name                                     string
		login, password, wantLogin, wantPassword string
		want                                     bool
	}{
		{name: "match", login: "admin", password: "new password", wantLogin: "admin", wantPassword: "new password", want: true},
		{name: "wrong login", login: "editor", password: "new password", wantLogin: "admin", wantPassword: "new password"},
		{name: "wrong password", login: "admin", password: "old password", wantLogin: "admin", wantPassword: "new password"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := validAdminCredentials(tt.login, tt.password, tt.wantLogin, tt.wantPassword); got != tt.want {
				t.Fatalf("validAdminCredentials() = %v, want %v", got, tt.want)
			}
		})
	}
}
