package static

import "embed"

//go:embed media/* dist/**/*
var Assets embed.FS

//go:embed widget/widget.js
var Widget []byte
