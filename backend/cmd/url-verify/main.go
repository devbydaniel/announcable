package main

import (
	"fmt"
	"github.com/devbydaniel/announcable/internal/util"
)

// Manual verification program for URL edge cases
// Run with: cd backend && go run internal/util/url_verification.go

func main() {
	fmt.Println("=== URL Utility Manual Verification ===\n")

	// Define test BASE_URL configurations
	testConfigs := []struct {
		name    string
		baseURL string
	}{
		{
			name:    "With trailing slash",
			baseURL: "https://example.com/",
		},
		{
			name:    "Without trailing slash",
			baseURL: "https://example.com",
		},
		{
			name:    "With port",
			baseURL: "https://example.com:8080",
		},
		{
			name:    "With subpath",
			baseURL: "https://example.com/app",
		},
		{
			name:    "Without protocol",
			baseURL: "example.com",
		},
		{
			name:    "Localhost with port (should use http)",
			baseURL: "localhost:3000",
		},
		{
			name:    "With port and subpath",
			baseURL: "https://example.com:8080/app",
		},
		{
			name:    "With port, subpath and trailing slash",
			baseURL: "https://example.com:8080/app/",
		},
	}

	// Test each configuration with real-world use cases
	for _, config := range testConfigs {
		fmt.Printf("Configuration: %s\n", config.name)
		fmt.Printf("BASE_URL: %s\n", config.baseURL)
		fmt.Println("---")

		// Test NormalizeBaseURL
		normalized, err := util.NormalizeBaseURL(config.baseURL)
		if err != nil {
			fmt.Printf("  ❌ NormalizeBaseURL ERROR: %v\n", err)
			fmt.Println()
			continue
		}
		fmt.Printf("  Normalized: %s\n", normalized)

		// Test real-world scenarios
		scenarios := []struct {
			name     string
			path     []string
			expected string
		}{
			{
				name:     "Email verification URL",
				path:     []string{"verify", "abc123token"},
				expected: normalized + "/verify/abc123token",
			},
			{
				name:     "Password reset URL",
				path:     []string{"reset-pw", "xyz789token"},
				expected: normalized + "/reset-pw/xyz789token",
			},
			{
				name:     "Invite accept URL",
				path:     []string{"invite-accept", "invite-token-123"},
				expected: normalized + "/invite-accept/invite-token-123",
			},
			{
				name:     "Release page URL",
				path:     []string{"s", "my-product-slug"},
				expected: normalized + "/s/my-product-slug",
			},
			{
				name:     "Image proxy URL",
				path:     []string{"api", "img"},
				expected: normalized + "/api/img",
			},
		}

		allPassed := true
		for _, scenario := range scenarios {
			result := util.BuildURL(config.baseURL, scenario.path...)
			passed := result == scenario.expected

			if passed {
				fmt.Printf("  ✓ %s: %s\n", scenario.name, result)
			} else {
				fmt.Printf("  ❌ %s\n", scenario.name)
				fmt.Printf("     Expected: %s\n", scenario.expected)
				fmt.Printf("     Got:      %s\n", result)
				allPassed = false
			}
		}

		// Check for common issues
		fmt.Println("\n  Edge case checks:")

		// Check for double slashes (except after protocol)
		testURL := util.BuildURL(config.baseURL, "api", "test")
		if containsDoubleSlash(testURL) {
			fmt.Printf("  ❌ Double slash detected in: %s\n", testURL)
			allPassed = false
		} else {
			fmt.Printf("  ✓ No double slashes\n")
		}

		// Check protocol is present
		if normalized[:4] != "http" {
			fmt.Printf("  ❌ Missing protocol in normalized URL\n")
			allPassed = false
		} else {
			fmt.Printf("  ✓ Protocol present\n")
		}

		// Check no trailing slash in normalized URL
		if len(normalized) > 0 && normalized[len(normalized)-1] == '/' {
			fmt.Printf("  ❌ Trailing slash in normalized URL\n")
			allPassed = false
		} else {
			fmt.Printf("  ✓ No trailing slash in normalized URL\n")
		}

		if allPassed {
			fmt.Printf("\n  ✅ All tests PASSED for this configuration\n")
		} else {
			fmt.Printf("\n  ❌ Some tests FAILED for this configuration\n")
		}

		fmt.Println("\n" + repeatString("=", 60) + "\n")
	}

	fmt.Println("=== Verification Complete ===")
}

// containsDoubleSlash checks for double slashes except after protocol
func containsDoubleSlash(url string) bool {
	// Remove protocol to avoid false positive
	if len(url) < 8 {
		return false
	}

	var checkPart string
	if url[:8] == "https://" {
		checkPart = url[8:]
	} else if url[:7] == "http://" {
		checkPart = url[7:]
	} else {
		checkPart = url
	}

	for i := 0; i < len(checkPart)-1; i++ {
		if checkPart[i] == '/' && checkPart[i+1] == '/' {
			return true
		}
	}
	return false
}

// repeatString repeats a string n times
func repeatString(s string, count int) string {
	result := ""
	for i := 0; i < count; i++ {
		result += s
	}
	return result
}
