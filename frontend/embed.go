package frontend

import "embed"

//go:embed dist/*
var staticFiles embed.FS

func GetStaticFiles() *embed.FS {
	return &staticFiles
}
