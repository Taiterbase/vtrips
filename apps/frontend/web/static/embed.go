package static

import "embed"

// Static contains the embedded static assets to be served by the application.
//
//go:embed * **/*
var Static embed.FS
