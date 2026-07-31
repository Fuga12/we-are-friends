package migrations

import _ "embed"

// Initial is compiled into the API image, so the container can initialise a
// fresh PostgreSQL database without a separate migration binary.
//
//go:embed 0001_init.sql
var Initial string
