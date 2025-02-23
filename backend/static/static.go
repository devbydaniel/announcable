package static

import "embed"

//go:embed css/**/* js/**/* media/*
var Assets embed.FS

//go:embed widget/widget.js
var Widget []byte
