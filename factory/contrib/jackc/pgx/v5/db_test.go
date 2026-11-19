package pgx

import "testing"

func TestNewConnConfig(t *testing.T) {
	t.Setenv("PGUSER", "root")
	t.Setenv("PGPASSWORD", "envpass")

	tests := []struct {
		name     string
		opts     Options
		wantUser string
		wantPass string
	}{
		{
			name:     "username and password override PGUSER and OS user when connect string has no user",
			opts:     Options{ConnectString: "postgres://db.example.com:5432/mydb?sslmode=disable", Username: "app", Password: "secret"},
			wantUser: "app",
			wantPass: "secret",
		},
		{
			name:     "username and password override the ones in the connect string",
			opts:     Options{ConnectString: "postgres://cs_user:cs_pass@db.example.com:5432/mydb?sslmode=disable", Username: "app", Password: "secret"},
			wantUser: "app",
			wantPass: "secret",
		},
		{
			name:     "empty username and password keep the connect string ones",
			opts:     Options{ConnectString: "postgres://cs_user:cs_pass@db.example.com:5432/mydb?sslmode=disable"},
			wantUser: "cs_user",
			wantPass: "cs_pass",
		},
		{
			name:     "only password set keeps the connect string user",
			opts:     Options{ConnectString: "postgres://cs_user:cs_pass@db.example.com:5432/mydb?sslmode=disable", Password: "secret"},
			wantUser: "cs_user",
			wantPass: "secret",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg, err := NewConnConfig(&tt.opts)
			if err != nil {
				t.Fatalf("NewConnConfig() error = %v", err)
			}
			if cfg.User != tt.wantUser {
				t.Errorf("User = %q, want %q", cfg.User, tt.wantUser)
			}
			if cfg.Password != tt.wantPass {
				t.Errorf("Password = %q, want %q", cfg.Password, tt.wantPass)
			}
		})
	}
}

func TestNewConnConfigInvalidConnectString(t *testing.T) {
	if _, err := NewConnConfig(&Options{ConnectString: "postgres://%zz"}); err == nil {
		t.Fatal("NewConnConfig() error = nil, want parse error")
	}
}
