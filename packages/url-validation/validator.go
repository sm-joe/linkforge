package urlvalidation

import (
	"context"
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
	ErrInvalidHostname   = errors.New("invalid hostname")
	ErrDNSResolution     = errors.New("DNS resolution failed")
)

type Resolver interface {
	LookupIP(
		ctx context.Context,
		network string,
		host string,
	) ([]net.IP, error)
}

type defaultResolver struct{}

func (defaultResolver) LookupIP(
	ctx context.Context,
	network string,
	host string,
) ([]net.IP, error) {
	return net.DefaultResolver.LookupIP(
		ctx,
		network,
		host,
	)
}

func Validate(rawURL string) error {
	rawURL = strings.TrimSpace(rawURL)

	if rawURL == "" {
		return ErrInvalidURL
	}

	parsed, err := url.ParseRequestURI(rawURL)
	if err != nil {
		return ErrInvalidURL
	}

	if parsed.Scheme != "http" &&
		parsed.Scheme != "https" {
		return ErrUnsupportedScheme
	}

	if parsed.Hostname() == "" {
		return ErrMissingHostname
	}

	hostname := strings.TrimSpace(
		parsed.Hostname(),
	)

	if !isValidHostname(hostname) {
		return ErrInvalidHostname
	}

	if ip := net.ParseIP(hostname); ip != nil {
		if isPrivateOrReservedIP(ip) {
			return ErrPrivateAddress
		}
	}

	return nil
}

func ValidateHTTPS(rawURL string) error {
	if err := Validate(rawURL); err != nil {
		return err
	}

	parsed, err := url.ParseRequestURI(
		strings.TrimSpace(rawURL),
	)
	if err != nil {
		return ErrInvalidURL
	}

	if parsed.Scheme != "https" {
		return ErrUnsupportedScheme
	}

	return nil
}

func ValidateContext(
	ctx context.Context,
	rawURL string,
) error {
	return validateContext(
		ctx,
		rawURL,
		defaultResolver{},
	)
}

func ValidateContextWithResolver(
	ctx context.Context,
	rawURL string,
	resolver Resolver,
) error {
	return validateContext(
		ctx,
		rawURL,
		resolver,
	)
}

func validateContext(
	ctx context.Context,
	rawURL string,
	resolver Resolver,
) error {
	if ctx == nil {
		ctx = context.Background()
	}

	if resolver == nil {
		resolver = defaultResolver{}
	}

	if err := Validate(rawURL); err != nil {
		return err
	}

	parsed, err := url.ParseRequestURI(
		strings.TrimSpace(rawURL),
	)
	if err != nil {
		return ErrInvalidURL
	}

	hostname := parsed.Hostname()

	ips, err := resolver.LookupIP(
		ctx,
		"ip",
		hostname,
	)
	if err != nil {
		return ErrDNSResolution
	}

	if len(ips) == 0 {
		return ErrDNSResolution
	}

	for _, ip := range ips {
		if isPrivateOrReservedIP(ip) {
			return ErrPrivateAddress
		}
	}

	return nil
}

func isValidHostname(hostname string) bool {
	if hostname == "" {
		return false
	}

	if len(hostname) > 253 {
		return false
	}

	if strings.HasSuffix(hostname, ".") {
		hostname = strings.TrimSuffix(
			hostname,
			".",
		)
	}

	if hostname == "" {
		return false
	}

	if net.ParseIP(hostname) != nil {
		return true
	}

	for _, label := range strings.Split(
		hostname,
		".",
	) {
		if !isValidHostnameLabel(label) {
			return false
		}
	}

	return true
}

func isValidHostnameLabel(label string) bool {
	if label == "" ||
		len(label) > 63 {
		return false
	}

	if label[0] == '-' ||
		label[len(label)-1] == '-' {
		return false
	}

	for _, char := range label {
		if (char >= 'a' && char <= 'z') ||
			(char >= 'A' && char <= 'Z') ||
			(char >= '0' && char <= '9') ||
			char == '-' {
			continue
		}

		return false
	}

	return true
}

var blockedIPv4Networks = mustParseCIDRs(
	"0.0.0.0/8",
	"10.0.0.0/8",
	"100.64.0.0/10",
	"127.0.0.0/8",
	"169.254.0.0/16",
	"172.16.0.0/12",
	"192.0.0.0/24",
	"192.0.2.0/24",
	"192.88.99.0/24",
	"192.168.0.0/16",
	"198.18.0.0/15",
	"198.51.100.0/24",
	"203.0.113.0/24",
	"224.0.0.0/4",
	"240.0.0.0/4",
)

var blockedIPv6Networks = mustParseCIDRs(
	"::/128",
	"::1/128",
	"100::/64",
	"2001:db8::/32",
	"fc00::/7",
	"fe80::/10",
	"ff00::/8",
)

func mustParseCIDRs(
	cidrs ...string,
) []*net.IPNet {
	networks := make(
		[]*net.IPNet,
		0,
		len(cidrs),
	)

	for _, cidr := range cidrs {
		_, network, err := net.ParseCIDR(cidr)
		if err != nil {
			panic(err)
		}

		networks = append(
			networks,
			network,
		)
	}

	return networks
}

func isPrivateOrReservedIP(ip net.IP) bool {
	if ip == nil {
		return true
	}

	if ipv4 := ip.To4(); ipv4 != nil {
		for _, network := range blockedIPv4Networks {
			if network.Contains(ipv4) {
				return true
			}
		}

		return false
	}

	ipv6 := ip.To16()
	if ipv6 == nil {
		return true
	}

	for _, network := range blockedIPv6Networks {
		if network.Contains(ipv6) {
			return true
		}
	}

	return false
}

func IsPrivateOrReservedIP(ip net.IP) bool {
	return isPrivateOrReservedIP(ip)
}
