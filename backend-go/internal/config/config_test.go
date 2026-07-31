package config

import "testing"

func TestLoadRejectsMissingRequiredValue(t *testing.T) {
	t.Setenv("DATABASE_URL", "")
	if _, err := Load(); err == nil {
		t.Fatal("Load() succeeded without DATABASE_URL")
	}
}

func TestLoadUsesDefaultsAndSplitsOrigins(t *testing.T) {
	for key, value := range map[string]string{
		"DATABASE_URL": "postgres://example", "MINIO_ENDPOINT": "localhost:9000",
		"MINIO_ACCESS_KEY": "access", "MINIO_SECRET_KEY": "secret", "MINIO_BUCKET": "media",
		"MEDIA_PUBLIC_BASE_URL": "http://localhost:9000/media", "ADMIN_USER": "admin", "ADMIN_PASSWORD": "password",
		"CORS_ORIGINS": "http://localhost:5173, https://example.test",
	} {
		t.Setenv(key, value)
	}
	cfg, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Port != "8080" || len(cfg.CORSOrigins) != 2 {
		t.Fatalf("unexpected config: %#v", cfg)
	}
}
