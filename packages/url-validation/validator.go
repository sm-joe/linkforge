package urlvalidation

import (
	"errors"
	"net"
	"net/url"
	"strings"
)

var (
	ErrInvalidURL        = errors.New("invalid URL")
	ErrUnsupportedScheme = errors.New("unsupported URL scheme")
	ErrMissingHostname   = errors.New("missing hostname")
	ErrPrivateAddress    = errors.New("private or reserved address")
)

func Validate(rawURL string) error {
	rawURL = strings.TrimSpace(rawURL)

	if rawURL == "" {
		return ErrInvalidURL
	}

	parsed, err := url.ParseRequestURI(rawURL)
	if err != nil {
		return ErrInvalidURL
	}

	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		return ErrUnsupportedScheme
	}

	if parsed.Hostname() == "" {
		return ErrMissingHostname
	}

	hostname := parsed.Hostname()

	if ip := net.ParseIP(hostname); ip != nil {
		if isPrivateOrReservedIP(ip) {
			return ErrPrivateAddress
		}
	}

	return nil
}

func isPrivateOrReservedIP(ip net.IP) bool {
	if ip.IsLoopback() ||
		ip.IsPrivate() ||
		ip.IsLinkLocalUnicast() ||
		ip.IsLinkLocalMulticast() ||
		ip.IsUnspecified() {
		return true
	}

	return false
}
