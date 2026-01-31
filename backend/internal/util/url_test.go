package util

import (
	"testing"
)

// TestNormalizeBaseURL_WithProtocol tests URLs that already have a protocol
func TestNormalizeBaseURL_WithProtocol(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
		wantErr  bool
	}{
		{
			name:     "https with domain",
			input:    "https://example.com",
			expected: "https://example.com",
			wantErr:  false,
		},
		{
			name:     "http with domain",
			input:    "http://example.com",
			expected: "http://example.com",
			wantErr:  false,
		},
		{
			name:     "https with trailing slash",
			input:    "https://example.com/",
			expected: "https://example.com",
			wantErr:  false,
		},
		{
			name:     "https with port",
			input:    "https://example.com:8080",
			expected: "https://example.com:8080",
			wantErr:  false,
		},
		{
			name:     "https with port and trailing slash",
			input:    "https://example.com:8080/",
			expected: "https://example.com:8080",
			wantErr:  false,
		},
		{
			name:     "https with subpath",
			input:    "https://example.com/subpath",
			expected: "https://example.com/subpath",
			wantErr:  false,
		},
		{
			name:     "https with subpath and trailing slash",
			input:    "https://example.com/subpath/",
			expected: "https://example.com/subpath",
			wantErr:  false,
		},
		{
			name:     "https with multiple subpaths",
			input:    "https://example.com/api/v1",
			expected: "https://example.com/api/v1",
			wantErr:  false,
		},
		{
			name:     "https with multiple subpaths and trailing slash",
			input:    "https://example.com/api/v1/",
			expected: "https://example.com/api/v1",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := NormalizeBaseURL(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("NormalizeBaseURL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if result != tt.expected {
				t.Errorf("NormalizeBaseURL() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

// TestNormalizeBaseURL_WithoutProtocol tests URLs without protocol (should add https by default)
func TestNormalizeBaseURL_WithoutProtocol(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
		wantErr  bool
	}{
		{
			name:     "domain without protocol",
			input:    "example.com",
			expected: "https://example.com",
			wantErr:  false,
		},
		{
			name:     "domain with trailing slash no protocol",
			input:    "example.com/",
			expected: "https://example.com",
			wantErr:  false,
		},
		{
			name:     "domain with port no protocol",
			input:    "example.com:8080",
			expected: "https://example.com:8080",
			wantErr:  false,
		},
		{
			name:     "domain with port and trailing slash no protocol",
			input:    "example.com:8080/",
			expected: "https://example.com:8080",
			wantErr:  false,
		},
		{
			name:     "domain with subpath no protocol",
			input:    "example.com/subpath",
			expected: "https://example.com/subpath",
			wantErr:  false,
		},
		{
			name:     "domain with subpath and trailing slash no protocol",
			input:    "example.com/subpath/",
			expected: "https://example.com/subpath",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := NormalizeBaseURL(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("NormalizeBaseURL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if result != tt.expected {
				t.Errorf("NormalizeBaseURL() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

// TestNormalizeBaseURL_Localhost tests localhost special handling (should use http)
func TestNormalizeBaseURL_Localhost(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
		wantErr  bool
	}{
		{
			name:     "localhost without protocol",
			input:    "localhost",
			expected: "http://localhost",
			wantErr:  false,
		},
		{
			name:     "localhost with port",
			input:    "localhost:3000",
			expected: "http://localhost:3000",
			wantErr:  false,
		},
		{
			name:     "localhost with trailing slash",
			input:    "localhost/",
			expected: "http://localhost",
			wantErr:  false,
		},
		{
			name:     "localhost with port and trailing slash",
			input:    "localhost:3000/",
			expected: "http://localhost:3000",
			wantErr:  false,
		},
		{
			name:     "localhost with subpath",
			input:    "localhost/api",
			expected: "http://localhost/api",
			wantErr:  false,
		},
		{
			name:     "localhost with explicit https",
			input:    "https://localhost",
			expected: "https://localhost",
			wantErr:  false,
		},
		{
			name:     "127.0.0.1 without protocol",
			input:    "127.0.0.1",
			expected: "http://127.0.0.1",
			wantErr:  false,
		},
		{
			name:     "127.0.0.1 with port",
			input:    "127.0.0.1:8080",
			expected: "http://127.0.0.1:8080",
			wantErr:  false,
		},
		{
			name:     "127.0.0.1 with trailing slash",
			input:    "127.0.0.1/",
			expected: "http://127.0.0.1",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := NormalizeBaseURL(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("NormalizeBaseURL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if result != tt.expected {
				t.Errorf("NormalizeBaseURL() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

// TestNormalizeBaseURL_ErrorCases tests invalid inputs that should return errors
func TestNormalizeBaseURL_ErrorCases(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{
			name:    "empty string",
			input:   "",
			wantErr: true,
		},
		{
			name:    "invalid URL characters",
			input:   "https://invalid url with spaces.com",
			wantErr: true,
		},
		{
			name:    "protocol only",
			input:   "https://",
			wantErr: true,
		},
		{
			name:    "invalid protocol",
			input:   "ftp://example.com",
			wantErr: false, // This should actually work - it's a valid URL
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NormalizeBaseURL(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("NormalizeBaseURL() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestBuildURL_BasicUsage tests basic URL building functionality
func TestBuildURL_BasicUsage(t *testing.T) {
	tests := []struct {
		name         string
		baseURL      string
		pathSegments []string
		expected     string
	}{
		{
			name:         "base with single segment",
			baseURL:      "https://example.com",
			pathSegments: []string{"api"},
			expected:     "https://example.com/api",
		},
		{
			name:         "base with multiple segments",
			baseURL:      "https://example.com",
			pathSegments: []string{"api", "v1", "users"},
			expected:     "https://example.com/api/v1/users",
		},
		{
			name:         "base with trailing slash and single segment",
			baseURL:      "https://example.com/",
			pathSegments: []string{"api"},
			expected:     "https://example.com/api",
		},
		{
			name:         "base with trailing slash and multiple segments",
			baseURL:      "https://example.com/",
			pathSegments: []string{"api", "v1", "users"},
			expected:     "https://example.com/api/v1/users",
		},
		{
			name:         "no path segments",
			baseURL:      "https://example.com",
			pathSegments: []string{},
			expected:     "https://example.com",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := BuildURL(tt.baseURL, tt.pathSegments...)
			if result != tt.expected {
				t.Errorf("BuildURL() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

// TestBuildURL_WithoutProtocol tests URL building when base lacks protocol
func TestBuildURL_WithoutProtocol(t *testing.T) {
	tests := []struct {
		name         string
		baseURL      string
		pathSegments []string
		expected     string
	}{
		{
			name:         "domain without protocol single segment",
			baseURL:      "example.com",
			pathSegments: []string{"api"},
			expected:     "https://example.com/api",
		},
		{
			name:         "domain without protocol multiple segments",
			baseURL:      "example.com",
			pathSegments: []string{"api", "v1", "users"},
			expected:     "https://example.com/api/v1/users",
		},
		{
			name:         "localhost without protocol",
			baseURL:      "localhost:3000",
			pathSegments: []string{"api"},
			expected:     "http://localhost:3000/api",
		},
		{
			name:         "127.0.0.1 without protocol",
			baseURL:      "127.0.0.1:8080",
			pathSegments: []string{"health"},
			expected:     "http://127.0.0.1:8080/health",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := BuildURL(tt.baseURL, tt.pathSegments...)
			if result != tt.expected {
				t.Errorf("BuildURL() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

// TestBuildURL_WithPort tests URL building with port numbers
func TestBuildURL_WithPort(t *testing.T) {
	tests := []struct {
		name         string
		baseURL      string
		pathSegments []string
		expected     string
	}{
		{
			name:         "https with standard port",
			baseURL:      "https://example.com:443",
			pathSegments: []string{"api"},
			expected:     "https://example.com:443/api",
		},
		{
			name:         "http with standard port",
			baseURL:      "http://example.com:80",
			pathSegments: []string{"api"},
			expected:     "http://example.com:80/api",
		},
		{
			name:         "custom port",
			baseURL:      "https://example.com:8080",
			pathSegments: []string{"api", "v1"},
			expected:     "https://example.com:8080/api/v1",
		},
		{
			name:         "custom port with trailing slash",
			baseURL:      "https://example.com:8080/",
			pathSegments: []string{"api"},
			expected:     "https://example.com:8080/api",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := BuildURL(tt.baseURL, tt.pathSegments...)
			if result != tt.expected {
				t.Errorf("BuildURL() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

// TestBuildURL_WithSubpath tests URL building when base has a subpath
func TestBuildURL_WithSubpath(t *testing.T) {
	tests := []struct {
		name         string
		baseURL      string
		pathSegments []string
		expected     string
	}{
		{
			name:         "base with subpath single segment",
			baseURL:      "https://example.com/app",
			pathSegments: []string{"api"},
			expected:     "https://example.com/app/api",
		},
		{
			name:         "base with subpath multiple segments",
			baseURL:      "https://example.com/app",
			pathSegments: []string{"api", "v1", "users"},
			expected:     "https://example.com/app/api/v1/users",
		},
		{
			name:         "base with subpath and trailing slash",
			baseURL:      "https://example.com/app/",
			pathSegments: []string{"api"},
			expected:     "https://example.com/app/api",
		},
		{
			name:         "base with nested subpath",
			baseURL:      "https://example.com/api/v1",
			pathSegments: []string{"users"},
			expected:     "https://example.com/api/v1/users",
		},
		{
			name:         "base with nested subpath and trailing slash",
			baseURL:      "https://example.com/api/v1/",
			pathSegments: []string{"users", "123"},
			expected:     "https://example.com/api/v1/users/123",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := BuildURL(tt.baseURL, tt.pathSegments...)
			if result != tt.expected {
				t.Errorf("BuildURL() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

// TestBuildURL_LeadingSlashes tests handling of leading slashes in path segments
func TestBuildURL_LeadingSlashes(t *testing.T) {
	tests := []struct {
		name         string
		baseURL      string
		pathSegments []string
		expected     string
	}{
		{
			name:         "segment with leading slash",
			baseURL:      "https://example.com",
			pathSegments: []string{"/api"},
			expected:     "https://example.com/api",
		},
		{
			name:         "multiple segments with leading slashes",
			baseURL:      "https://example.com",
			pathSegments: []string{"/api", "/v1", "/users"},
			expected:     "https://example.com/api/v1/users",
		},
		{
			name:         "mixed segments some with leading slashes",
			baseURL:      "https://example.com",
			pathSegments: []string{"/api", "v1", "/users"},
			expected:     "https://example.com/api/v1/users",
		},
		{
			name:         "base with trailing slash and segment with leading slash",
			baseURL:      "https://example.com/",
			pathSegments: []string{"/api"},
			expected:     "https://example.com/api",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := BuildURL(tt.baseURL, tt.pathSegments...)
			if result != tt.expected {
				t.Errorf("BuildURL() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

// TestBuildURL_EmptySegments tests handling of empty path segments
func TestBuildURL_EmptySegments(t *testing.T) {
	tests := []struct {
		name         string
		baseURL      string
		pathSegments []string
		expected     string
	}{
		{
			name:         "single empty segment",
			baseURL:      "https://example.com",
			pathSegments: []string{""},
			expected:     "https://example.com",
		},
		{
			name:         "empty segments between valid ones",
			baseURL:      "https://example.com",
			pathSegments: []string{"api", "", "users"},
			expected:     "https://example.com/api/users",
		},
		{
			name:         "multiple empty segments",
			baseURL:      "https://example.com",
			pathSegments: []string{"", "", ""},
			expected:     "https://example.com",
		},
		{
			name:         "empty segments at start and end",
			baseURL:      "https://example.com",
			pathSegments: []string{"", "api", ""},
			expected:     "https://example.com/api",
		},
		{
			name:         "whitespace-only segments",
			baseURL:      "https://example.com",
			pathSegments: []string{"api", "  ", "users"},
			expected:     "https://example.com/api/users",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := BuildURL(tt.baseURL, tt.pathSegments...)
			if result != tt.expected {
				t.Errorf("BuildURL() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

// TestBuildURL_EdgeCases tests edge cases and error conditions
func TestBuildURL_EdgeCases(t *testing.T) {
	tests := []struct {
		name         string
		baseURL      string
		pathSegments []string
		expected     string
	}{
		{
			name:         "empty base URL",
			baseURL:      "",
			pathSegments: []string{"api"},
			expected:     "",
		},
		{
			name:         "empty base and empty segments",
			baseURL:      "",
			pathSegments: []string{"", ""},
			expected:     "",
		},
		{
			name:         "token as path segment",
			baseURL:      "https://example.com",
			pathSegments: []string{"verify-email", "abc123token"},
			expected:     "https://example.com/verify-email/abc123token",
		},
		{
			name:         "UUID as path segment",
			baseURL:      "https://example.com",
			pathSegments: []string{"users", "550e8400-e29b-41d4-a716-446655440000"},
			expected:     "https://example.com/users/550e8400-e29b-41d4-a716-446655440000",
		},
		{
			name:         "slug as path segment",
			baseURL:      "https://example.com",
			pathSegments: []string{"s", "my-release-page"},
			expected:     "https://example.com/s/my-release-page",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := BuildURL(tt.baseURL, tt.pathSegments...)
			if result != tt.expected {
				t.Errorf("BuildURL() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

// TestBuildURL_RealWorldScenarios tests actual use cases from the application
func TestBuildURL_RealWorldScenarios(t *testing.T) {
	tests := []struct {
		name         string
		baseURL      string
		pathSegments []string
		expected     string
		description  string
	}{
		{
			name:         "email verification URL",
			baseURL:      "https://app.example.com",
			pathSegments: []string{"verify-email", "token123"},
			expected:     "https://app.example.com/verify-email/token123",
			description:  "User service email verification",
		},
		{
			name:         "password reset URL",
			baseURL:      "https://app.example.com",
			pathSegments: []string{"reset-pw", "resettoken456"},
			expected:     "https://app.example.com/reset-pw/resettoken456",
			description:  "Password reset handler",
		},
		{
			name:         "invite accept URL",
			baseURL:      "https://app.example.com",
			pathSegments: []string{"invite-accept", "invitetoken789"},
			expected:     "https://app.example.com/invite-accept/invitetoken789",
			description:  "Organisation invite service",
		},
		{
			name:         "release page URL",
			baseURL:      "https://app.example.com",
			pathSegments: []string{"s", "my-product-updates"},
			expected:     "https://app.example.com/s/my-product-updates",
			description:  "Release page config service",
		},
		{
			name:         "image proxy URL",
			baseURL:      "https://app.example.com",
			pathSegments: []string{"api", "img"},
			expected:     "https://app.example.com/api/img",
			description:  "Object storage service image proxy",
		},
		{
			name:         "self-hosted with trailing slash",
			baseURL:      "https://my-company.com/announcable/",
			pathSegments: []string{"verify-email", "token"},
			expected:     "https://my-company.com/announcable/verify-email/token",
			description:  "Self-hosted deployment with subpath",
		},
		{
			name:         "local development",
			baseURL:      "localhost:3000",
			pathSegments: []string{"api", "v1", "health"},
			expected:     "http://localhost:3000/api/v1/health",
			description:  "Local development server",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := BuildURL(tt.baseURL, tt.pathSegments...)
			if result != tt.expected {
				t.Errorf("BuildURL() = %v, expected %v (scenario: %s)", result, tt.expected, tt.description)
			}
		})
	}
}
