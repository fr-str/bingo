package migrations

import "embed"

//go:embed migrations
var DBmigrations embed.FS
