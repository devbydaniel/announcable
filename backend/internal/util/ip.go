package util

import (
	"net"
	"net/http"
	"strings"
)

// GetClientIP extracts the client's IP address from the HTTP request.
// It checks headers in the following order for proxy support:
// 1. X-Forwarded-For (takes the first IP if multiple are present)
// 2. X-Real-IP
// 3. RemoteAddr (with port stripped)
// Handles both IPv4 and IPv6 addresses correctly.
func GetClientIP(r *http.Request) string {
	// Check X-Forwarded-For header first (for requests through proxies/load balancers)
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		// X-Forwarded-For can contain multiple IPs (client, proxy1, proxy2, ...)
		// We want the first one (the original client)
		ips := strings.Split(xff, ",")
		if len(ips) > 0 {
			ip := strings.TrimSpace(ips[0])
			if ip != "" {
				return ip
			}
		}
	}

	// Check X-Real-IP header (set by some proxies)
	if xri := r.Header.Get("X-Real-IP"); xri != "" {
		ip := strings.TrimSpace(xri)
		if ip != "" {
			return ip
		}
	}

	// Fall back to RemoteAddr
	// RemoteAddr is in the format "IP:port" or "[IPv6]:port"
	// We need to strip the port
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		// If SplitHostPort fails, it might be just an IP without port
		// This can happen in some test scenarios
		return r.RemoteAddr
	}

	return ip
}
