package fileEmbed

import (
	"embed"
)

//go:embed plugin
var Plugin embed.FS

//go:embed template
var Template embed.FS
