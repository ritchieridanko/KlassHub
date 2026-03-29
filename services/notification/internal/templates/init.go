package templates

import "embed"

//go:embed *.html.tmpl
var Email embed.FS
