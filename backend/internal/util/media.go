package util

import "strings"

// TransformMediaLink converts a YouTube or Loom share URL into its corresponding embed URL
func TransformMediaLink(mediaLink string) string {
	if mediaLink == "" {
		return ""
	}

	if strings.Contains(mediaLink, "youtube.com") {
		parts := strings.Split(mediaLink, "v=")
		if len(parts) < 2 {
			return ""
		}
		videoID := strings.Split(parts[1], "&")[0]
		return "https://www.youtube.com/embed/" + videoID
	}

	if strings.Contains(mediaLink, "loom.com") {
		return strings.Replace(mediaLink, "share", "embed", 1)
	}

	return ""
}
