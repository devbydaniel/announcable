package static

import "embed"

//go:embed css/**/* js/**/*
var Assets embed.FS

//go:embed widget/widget.js
var Widget []byte
