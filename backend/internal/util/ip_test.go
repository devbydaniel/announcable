package util

import (
	"net/http"
	"testing"
)

func TestGetClientIP(t *testing.T) {
	tests := []struct {
		name           string
		xForwardedFor  string
		xRealIP        string
		remoteAddr     string
		expectedIP     string
	}{
		{
			name:           "X-Forwarded-For with single IP",
			xForwardedFor:  "203.0.113.1",
			xRealIP:        "",
			remoteAddr:     "192.168.1.1:12345",
			expectedIP:     "203.0.113.1",
		},
		{
			name:           "X-Forwarded-For with multiple IPs (returns first)",
			xForwardedFor:  "203.0.113.1, 198.51.100.1, 192.0.2.1",
			xRealIP:        "",
			remoteAddr:     "192.168.1.1:12345",
			expectedIP:     "203.0.113.1",
		},
		{
			name:           "X-Forwarded-For with spaces around IPs",
			xForwardedFor:  "  203.0.113.1  , 198.51.100.1",
			xRealIP:        "",
			remoteAddr:     "192.168.1.1:12345",
			expectedIP:     "203.0.113.1",
		},
		{
			name:           "X-Real-IP when X-Forwarded-For is absent",
			xForwardedFor:  "",
			xRealIP:        "203.0.113.1",
			remoteAddr:     "192.168.1.1:12345",
			expectedIP:     "203.0.113.1",
		},
		{
			name:           "X-Forwarded-For takes precedence over X-Real-IP",
			xForwardedFor:  "203.0.113.1",
			xRealIP:        "198.51.100.1",
			remoteAddr:     "192.168.1.1:12345",
			expectedIP:     "203.0.113.1",
		},
		{
			name:           "RemoteAddr IPv4 with port when no headers present",
			xForwardedFor:  "",
			xRealIP:        "",
			remoteAddr:     "203.0.113.1:12345",
			expectedIP:     "203.0.113.1",
		},
		{
			name:           "RemoteAddr IPv6 with port when no headers present",
			xForwardedFor:  "",
			xRealIP:        "",
			remoteAddr:     "[2001:db8::1]:12345",
			expectedIP:     "2001:db8::1",
		},
		{
			name:           "RemoteAddr without port (edge case)",
			xForwardedFor:  "",
			xRealIP:        "",
			remoteAddr:     "203.0.113.1",
			expectedIP:     "203.0.113.1",
		},
		{
			name:           "Empty X-Forwarded-For falls back to X-Real-IP",
			xForwardedFor:  "",
			xRealIP:        "203.0.113.1",
			remoteAddr:     "192.168.1.1:12345",
			expectedIP:     "203.0.113.1",
		},
		{
			name:           "Whitespace-only X-Forwarded-For falls back to X-Real-IP",
			xForwardedFor:  "   ",
			xRealIP:        "203.0.113.1",
			remoteAddr:     "192.168.1.1:12345",
			expectedIP:     "203.0.113.1",
		},
		{
			name:           "IPv6 in X-Forwarded-For",
			xForwardedFor:  "2001:db8::1",
			xRealIP:        "",
			remoteAddr:     "192.168.1.1:12345",
			expectedIP:     "2001:db8::1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a mock HTTP request
			req, err := http.NewRequest("GET", "http://example.com", nil)
			if err != nil {
				t.Fatalf("Failed to create request: %v", err)
			}

			// Set headers if provided
			if tt.xForwardedFor != "" {
				req.Header.Set("X-Forwarded-For", tt.xForwardedFor)
			}
			if tt.xRealIP != "" {
				req.Header.Set("X-Real-IP", tt.xRealIP)
			}
			req.RemoteAddr = tt.remoteAddr

			// Call the function
			result := GetClientIP(req)

			// Verify the result
			if result != tt.expectedIP {
				t.Errorf("GetClientIP() = %v, want %v", result, tt.expectedIP)
			}
		})
	}
}
