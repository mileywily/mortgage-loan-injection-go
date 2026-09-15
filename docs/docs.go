package docs

import "embed"

//go:embed swagger.json swagger-ui.html
var Files embed.FS
