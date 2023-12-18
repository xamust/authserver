package www

import "embed"

//go:embed html
var HTML embed.FS

//go:embed scripts
var Scripts embed.FS

//go:embed styles
var Style embed.FS
