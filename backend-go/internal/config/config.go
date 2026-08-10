package config

import (
	"fmt"
	"net/url"
	"os"
	"strings"
)

// Config contains all process-level settings. Secrets are never logged.
type Config struct {
	Environment        string
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
	environment := valueOr("APP_ENV", "development")
	defaults, ok := profiles[environment]
	if !ok {
		return Config{}, fmt.Errorf("APP_ENV must be development or production")
	}
	databaseURL, err := databaseURL(defaults)
	if err != nil {
		return Config{}, err
	}
	cfg := Config{
		Environment:        environment,
		DatabaseURL:        databaseURL,
		MinIOEndpoint:      valueOr("MINIO_ENDPOINT", defaults.MinIOEndpoint),
		MinIOAccessKey:     os.Getenv("MINIO_ACCESS_KEY"),
		MinIOSecretKey:     os.Getenv("MINIO_SECRET_KEY"),
		MinIOBucket:        valueOr("MINIO_BUCKET", defaults.MinIOBucket),
		MediaPublicBaseURL: valueOr("MEDIA_PUBLIC_BASE_URL", defaults.MediaPublicBaseURL),
		Port:               valueOr("PORT", defaults.Port),
		CORSOrigins:        splitCSV(valueOr("CORS_ORIGINS", defaults.CORSOrigins)),
		AdminUser:          valueOr("ADMIN_USER", defaults.AdminUser),
		AdminPassword:      os.Getenv("ADMIN_PASSWORD"),
	}

	for name, value := range map[string]string{
		"DATABASE_URL":     cfg.DatabaseURL,
		"MINIO_ENDPOINT":   cfg.MinIOEndpoint,
		"MINIO_ACCESS_KEY": cfg.MinIOAccessKey,
		"MINIO_SECRET_KEY": cfg.MinIOSecretKey,
		"ADMIN_PASSWORD":   cfg.AdminPassword,
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

func databaseURL(defaults profile) (string, error) {
	if value := strings.TrimSpace(os.Getenv("DATABASE_URL")); value != "" {
		return value, nil
	}
	if defaults.DatabaseURLPattern == "" {
		return defaults.DatabaseURL, nil
	}
	password := os.Getenv("POSTGRES_PASSWORD")
	if password == "" {
		return "", fmt.Errorf("POSTGRES_PASSWORD must be set")
	}
	return fmt.Sprintf(defaults.DatabaseURLPattern, url.QueryEscape(password)), nil
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
