package util

import (
	"errors"
	"net/url"
	"strings"
)

// NormalizeBaseURL ensures the base URL has a protocol, removes trailing slashes, and validates format
func NormalizeBaseURL(baseURL string) (string, error) {
	if baseURL == "" {
		return "", errors.New("base URL cannot be empty")
	}

	// Add protocol if missing
	if !strings.HasPrefix(baseURL, "http://") && !strings.HasPrefix(baseURL, "https://") {
		// Default to https for production, but preserve localhost behavior
		if strings.HasPrefix(baseURL, "localhost") || strings.HasPrefix(baseURL, "127.0.0.1") {
			baseURL = "http://" + baseURL
		} else {
			baseURL = "https://" + baseURL
		}
	}

	// Parse to validate URL format
	parsedURL, err := url.Parse(baseURL)
	if err != nil {
		return "", errors.New("invalid base URL format: " + err.Error())
	}

	if parsedURL.Host == "" {
		return "", errors.New("base URL must include a host")
	}

	// Remove trailing slash from path
	parsedURL.Path = strings.TrimSuffix(parsedURL.Path, "/")

	return parsedURL.String(), nil
}

// BuildURL constructs a complete URL by joining the base URL with path segments
// It handles trailing slashes in base, leading slashes in segments, and empty segments
func BuildURL(baseURL string, pathSegments ...string) string {
	if baseURL == "" {
		return ""
	}

	// Normalize the base URL
	normalizedBase, err := NormalizeBaseURL(baseURL)
	if err != nil {
		// If normalization fails, try to work with what we have
		normalizedBase = strings.TrimSuffix(baseURL, "/")
	}

	// Filter out empty segments and clean leading slashes
	var cleanSegments []string
	for _, segment := range pathSegments {
		cleaned := strings.TrimSpace(segment)
		if cleaned != "" {
			// Remove leading slash from segment
			cleaned = strings.TrimPrefix(cleaned, "/")
			cleanSegments = append(cleanSegments, cleaned)
		}
	}

	// If no segments, return the base URL
	if len(cleanSegments) == 0 {
		return normalizedBase
	}

	// Join all segments with forward slashes
	path := strings.Join(cleanSegments, "/")

	// Combine base URL and path
	return normalizedBase + "/" + path
}
