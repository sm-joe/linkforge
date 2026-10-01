package link

import "strings"

func DetectBrowser(userAgent string) string {
	ua := strings.ToLower(userAgent)

	switch {
	case strings.Contains(ua, "edg/"):
		return "Edge"

	case strings.Contains(ua, "opr/") ||
		strings.Contains(ua, "opera"):
		return "Opera"

	case strings.Contains(ua, "chrome/") &&
		!strings.Contains(ua, "edg/"):
		return "Chrome"

	case strings.Contains(ua, "firefox/"):
		return "Firefox"

	case strings.Contains(ua, "safari/") &&
		!strings.Contains(ua, "chrome/"):
		return "Safari"

	default:
		return "Other"
	}
}

func DetectDevice(userAgent string) string {
	ua := strings.ToLower(userAgent)

	switch {
	case strings.Contains(ua, "ipad"):
		return "Tablet"

	case strings.Contains(ua, "tablet"):
		return "Tablet"

	case strings.Contains(ua, "iphone"):
		return "Mobile"

	case strings.Contains(ua, "android"):
		return "Mobile"

	case strings.Contains(ua, "mobile"):
		return "Mobile"

	default:
		return "Desktop"
	}
}