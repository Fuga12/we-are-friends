package config

import "testing"

func TestLoadRejectsMissingRequiredValue(t *testing.T) {
	t.Setenv("MINIO_ACCESS_KEY", "")
	if _, err := Load(); err == nil {
		t.Fatal("Load() succeeded without MINIO_ACCESS_KEY")
	}
}

func TestLoadUsesDefaultsAndSplitsOrigins(t *testing.T) {
	for key, value := range map[string]string{
		"MINIO_ACCESS_KEY": "access", "MINIO_SECRET_KEY": "secret", "ADMIN_PASSWORD": "password",
		"CORS_ORIGINS": "http://localhost:5173, https://example.test",
	} {
		t.Setenv(key, value)
	}
	cfg, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Environment != "development" || cfg.Port != "8080" || len(cfg.CORSOrigins) != 2 {
		t.Fatalf("unexpected config: %#v", cfg)
	}
}

func TestLoadProductionProfile(t *testing.T) {
	for key, value := range map[string]string{
		"APP_ENV": "production", "MINIO_ACCESS_KEY": "access",
		"MINIO_SECRET_KEY": "secret", "ADMIN_PASSWORD": "password", "POSTGRES_PASSWORD": "db secret",
	} {
		t.Setenv(key, value)
	}
	cfg, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	if cfg.MinIOEndpoint != "minio:9000" || cfg.MediaPublicBaseURL != "/media" ||
		cfg.DatabaseURL != "postgres://we_are_friends:db+secret@postgres:5432/we_are_friends?sslmode=disable" {
		t.Fatalf("unexpected production config: %#v", cfg)
	}
}
