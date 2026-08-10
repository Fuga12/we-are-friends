package config

var productionProfile = profile{
	DatabaseURLPattern: "postgres://we_are_friends:%s@postgres:5432/we_are_friends?sslmode=disable",
	MinIOEndpoint:      "minio:9000",
	MinIOBucket:        "media",
	MediaPublicBaseURL: "/media",
	Port:               "8080",
	CORSOrigins:        "http://localhost",
	AdminUser:          "admin",
}
