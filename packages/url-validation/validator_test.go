package urlvalidation

import (
	"context"
	"errors"
	"net"
	"testing"
)

type fakeResolver struct {
	ips []net.IP
	err error
}

func (r fakeResolver) LookupIP(
	ctx context.Context,
	network string,
	host string,
) ([]net.IP, error) {
	return r.ips, r.err
}

func TestValidate(t *testing.T) {
	tests := []struct {
		name    string
		rawURL  string
		wantErr error
	}{
		{
			name:   "valid HTTPS URL",
			rawURL: "https://example.com",
		},
		{
			name:   "valid HTTP URL",
			rawURL: "http://example.com",
		},
		{
			name:    "empty URL",
			rawURL:  "",
			wantErr: ErrInvalidURL,
		},
		{
			name:    "unsupported scheme",
			rawURL:  "ftp://example.com",
			wantErr: ErrUnsupportedScheme,
		},
		{
			name:    "missing hostname",
			rawURL:  "https:///path",
			wantErr: ErrMissingHostname,
		},
		{
			name:    "localhost",
			rawURL:  "http://127.0.0.1",
			wantErr: ErrPrivateAddress,
		},
		{
			name:    "private IPv4",
			rawURL:  "http://192.168.1.10",
			wantErr: ErrPrivateAddress,
		},
		{
			name:    "private IPv4 10 range",
			rawURL:  "http://10.0.0.10",
			wantErr: ErrPrivateAddress,
		},
		{
			name:    "private IPv4 172 range",
			rawURL:  "http://172.16.0.10",
			wantErr: ErrPrivateAddress,
		},
		{
			name:    "IPv6 loopback",
			rawURL:  "http://[::1]",
			wantErr: ErrPrivateAddress,
		},
		{
			name:    "invalid hostname with underscore",
			rawURL:  "https://bad_host.example.com",
			wantErr: ErrInvalidHostname,
		},
		{
			name:    "invalid hostname with leading hyphen",
			rawURL:  "https://-example.com",
			wantErr: ErrInvalidHostname,
		},
		{
			name:    "invalid hostname with trailing hyphen",
			rawURL:  "https://example-.com",
			wantErr: ErrInvalidHostname,
		},
		{
			name:    "invalid hostname with empty label",
			rawURL:  "https://example..com",
			wantErr: ErrInvalidHostname,
		},
		{
			name:   "valid subdomain",
			rawURL: "https://app.example.com",
		},
		{
			name:   "valid hostname with hyphen",
			rawURL: "https://my-example.com",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := Validate(test.rawURL)

			if test.wantErr == nil {
				if err != nil {
					t.Fatalf(
						"expected no error, got %v",
						err,
					)
				}

				return
			}

			if !errors.Is(err, test.wantErr) {
				t.Fatalf(
					"expected %v, got %v",
					test.wantErr,
					err,
				)
			}
		})
	}
}

func TestValidateHTTPS(t *testing.T) {
	tests := []struct {
		name    string
		rawURL  string
		wantErr error
	}{
		{
			name:   "valid HTTPS URL",
			rawURL: "https://example.com",
		},
		{
			name:    "HTTP URL rejected",
			rawURL:  "http://example.com",
			wantErr: ErrUnsupportedScheme,
		},
		{
			name:    "private HTTPS address rejected",
			rawURL:  "https://127.0.0.1",
			wantErr: ErrPrivateAddress,
		},
		{
			name:    "invalid HTTPS URL",
			rawURL:  "https://bad_host.example.com",
			wantErr: ErrInvalidHostname,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := ValidateHTTPS(test.rawURL)

			if test.wantErr == nil {
				if err != nil {
					t.Fatalf(
						"expected no error, got %v",
						err,
					)
				}

				return
			}

			if !errors.Is(err, test.wantErr) {
				t.Fatalf(
					"expected %v, got %v",
					test.wantErr,
					err,
				)
			}
		})
	}
}

func TestValidateContextWithResolver(
	t *testing.T,
) {
	tests := []struct {
		name     string
		rawURL   string
		resolver fakeResolver
		wantErr  error
	}{
		{
			name:   "public DNS address",
			rawURL: "https://example.com",
			resolver: fakeResolver{
				ips: []net.IP{
					net.ParseIP("93.184.216.34"),
				},
			},
		},
		{
			name:   "DNS resolves to private IPv4",
			rawURL: "https://example.com",
			resolver: fakeResolver{
				ips: []net.IP{
					net.ParseIP("10.0.0.10"),
				},
			},
			wantErr: ErrPrivateAddress,
		},
		{
			name:   "DNS resolves to loopback",
			rawURL: "https://example.com",
			resolver: fakeResolver{
				ips: []net.IP{
					net.ParseIP("127.0.0.1"),
				},
			},
			wantErr: ErrPrivateAddress,
		},
		{
			name:   "DNS resolves to IPv6 loopback",
			rawURL: "https://example.com",
			resolver: fakeResolver{
				ips: []net.IP{
					net.ParseIP("::1"),
				},
			},
			wantErr: ErrPrivateAddress,
		},
		{
			name:   "one public and one private address",
			rawURL: "https://example.com",
			resolver: fakeResolver{
				ips: []net.IP{
					net.ParseIP("93.184.216.34"),
					net.ParseIP("192.168.1.10"),
				},
			},
			wantErr: ErrPrivateAddress,
		},
		{
			name:   "DNS resolution failure",
			rawURL: "https://example.com",
			resolver: fakeResolver{
				err: errors.New("lookup failed"),
			},
			wantErr: ErrDNSResolution,
		},
		{
			name:   "DNS returned no addresses",
			rawURL: "https://example.com",
			resolver: fakeResolver{
				ips: []net.IP{},
			},
			wantErr: ErrDNSResolution,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := ValidateContextWithResolver(
				context.Background(),
				test.rawURL,
				test.resolver,
			)

			if test.wantErr == nil {
				if err != nil {
					t.Fatalf(
						"expected no error, got %v",
						err,
					)
				}

				return
			}

			if !errors.Is(err, test.wantErr) {
				t.Fatalf(
					"expected %v, got %v",
					test.wantErr,
					err,
				)
			}
		})
	}
}

func TestIsPrivateOrReservedIP(t *testing.T) {
	tests := []struct {
		name string
		ip   string
		want bool
	}{
		{
			name: "IPv4 unspecified",
			ip:   "0.0.0.0",
			want: true,
		},
		{
			name: "IPv4 loopback",
			ip:   "127.0.0.1",
			want: true,
		},
		{
			name: "IPv4 private 10",
			ip:   "10.10.10.10",
			want: true,
		},
		{
			name: "IPv4 private 172",
			ip:   "172.16.10.10",
			want: true,
		},
		{
			name: "IPv4 private 192",
			ip:   "192.168.10.10",
			want: true,
		},
		{
			name: "IPv4 CGNAT",
			ip:   "100.64.10.10",
			want: true,
		},
		{
			name: "IPv4 link local",
			ip:   "169.254.10.10",
			want: true,
		},
		{
			name: "IPv4 TEST-NET-1",
			ip:   "192.0.2.10",
			want: true,
		},
		{
			name: "IPv4 TEST-NET-2",
			ip:   "198.51.100.10",
			want: true,
		},
		{
			name: "IPv4 TEST-NET-3",
			ip:   "203.0.113.10",
			want: true,
		},
		{
			name: "IPv4 multicast",
			ip:   "224.0.0.1",
			want: true,
		},
		{
			name: "IPv4 reserved",
			ip:   "240.0.0.1",
			want: true,
		},
		{
			name: "IPv6 unspecified",
			ip:   "::",
			want: true,
		},
		{
			name: "IPv6 loopback",
			ip:   "::1",
			want: true,
		},
		{
			name: "IPv6 unique local",
			ip:   "fd00::10",
			want: true,
		},
		{
			name: "IPv6 link local",
			ip:   "fe80::10",
			want: true,
		},
		{
			name: "IPv6 documentation",
			ip:   "2001:db8::10",
			want: true,
		},
		{
			name: "IPv6 multicast",
			ip:   "ff02::1",
			want: true,
		},
		{
			name: "public IPv4",
			ip:   "8.8.8.8",
			want: false,
		},
		{
			name: "public IPv6",
			ip:   "2001:4860:4860::8888",
			want: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ip := net.ParseIP(test.ip)

			if ip == nil {
				t.Fatalf(
					"failed to parse test IP %q",
					test.ip,
				)
			}

			got := isPrivateOrReservedIP(ip)

			if got != test.want {
				t.Fatalf(
					"expected %v, got %v",
					test.want,
					got,
				)
			}
		})
	}
}
