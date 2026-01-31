package static

import "embed"

// Assets contains the embedded static files (media and dist).
//
//go:embed media/* dist/* dist/**/*
var Assets embed.FS

// Widget contains the embedded widget JavaScript bundle.
//
//go:embed widget/widget.js
var Widget []byte
