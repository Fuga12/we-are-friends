package config

import (
	"fmt"
	"os"
	"strings"
)

// Config contains all process-level settings. Secrets are never logged.
type Config struct {
	DatabaseURL        string
	MinIOEndpoint      string
	MinIOAccessKey     string
	MinIOSecretKey     string
	MinIOBucket        string
	MediaPublicBaseURL string
	Port               string
	CORSOrigins        []string
	AdminUser          string
	AdminPassword      string
}

// Load validates configuration at startup so an API instance never starts with
// an unintentionally open admin endpoint or an incomplete media configuration.
func Load() (Config, error) {
	cfg := Config{
		DatabaseURL:        os.Getenv("DATABASE_URL"),
		MinIOEndpoint:      os.Getenv("MINIO_ENDPOINT"),
		MinIOAccessKey:     os.Getenv("MINIO_ACCESS_KEY"),
		MinIOSecretKey:     os.Getenv("MINIO_SECRET_KEY"),
		MinIOBucket:        os.Getenv("MINIO_BUCKET"),
		MediaPublicBaseURL: os.Getenv("MEDIA_PUBLIC_BASE_URL"),
		Port:               valueOr("PORT", "8080"),
		CORSOrigins:        splitCSV(valueOr("CORS_ORIGINS", "http://localhost:5173")),
		AdminUser:          os.Getenv("ADMIN_USER"),
		AdminPassword:      os.Getenv("ADMIN_PASSWORD"),
	}

	for name, value := range map[string]string{
		"DATABASE_URL":          cfg.DatabaseURL,
		"MINIO_ENDPOINT":        cfg.MinIOEndpoint,
		"MINIO_ACCESS_KEY":      cfg.MinIOAccessKey,
		"MINIO_SECRET_KEY":      cfg.MinIOSecretKey,
		"MINIO_BUCKET":          cfg.MinIOBucket,
		"MEDIA_PUBLIC_BASE_URL": cfg.MediaPublicBaseURL,
		"ADMIN_USER":            cfg.AdminUser,
		"ADMIN_PASSWORD":        cfg.AdminPassword,
	} {
		if value == "" {
			return Config{}, fmt.Errorf("%s must be set", name)
		}
	}
	if len(cfg.CORSOrigins) == 0 {
		return Config{}, fmt.Errorf("CORS_ORIGINS must include at least one origin")
	}
	return cfg, nil
}

func valueOr(name, fallback string) string {
	if value := strings.TrimSpace(os.Getenv(name)); value != "" {
		return value
	}
	return fallback
}

func splitCSV(value string) []string {
	parts := strings.Split(value, ",")
	origins := make([]string, 0, len(parts))
	for _, part := range parts {
		if origin := strings.TrimSpace(part); origin != "" {
			origins = append(origins, origin)
		}
	}
	return origins
}
