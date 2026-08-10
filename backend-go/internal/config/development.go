package config

var developmentProfile = profile{
	DatabaseURL:        "postgres://we_are_friends:we_are_friends@localhost:15432/we_are_friends?sslmode=disable",
	MinIOEndpoint:      "localhost:9000",
	MinIOBucket:        "media",
	MediaPublicBaseURL: "http://localhost:9000/media",
	Port:               "8080",
	CORSOrigins:        "http://localhost:5173",
	AdminUser:          "admin",
}
