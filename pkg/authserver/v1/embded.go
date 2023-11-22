package authserver

import "embed"

//go:embed api.swagger.json
var SwaggerJsonApi embed.FS
