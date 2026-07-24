package web

import "embed"

// Files contains the production frontend built by Vite.
//
//go:embed dist/*
var Files embed.FS
