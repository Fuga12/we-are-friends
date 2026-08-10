package migrations

import _ "embed"

// SQL migrations are compiled into the API image, so the container can update
// an existing PostgreSQL database without a separate migration binary.
//
//go:embed 0001_init.sql
var Initial string

//go:embed 0002_admin_auth_and_publication.sql
var AdminAuthAndPublication string
