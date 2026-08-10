package config

// profile contains non-secret settings committed with the application.
// Every value can still be overridden by an environment variable in Load.
type profile struct {
	DatabaseURL        string
	DatabaseURLPattern string
	MinIOEndpoint      string
	MinIOBucket        string
	MediaPublicBaseURL string
	Port               string
	CORSOrigins        string
	AdminUser          string
}

var profiles = map[string]profile{
	"development": developmentProfile,
	"production":  productionProfile,
}
