package link

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"time"

	urlvalidation "github.com/sm-joe/linkforge/packages/url-validation"
)

const healthRequestTimeout = 10 * time.Second

type HealthClient struct {
	client   *http.Client
	resolver urlvalidation.Resolver
	validate bool
}

func NewHealthClient(
	resolver urlvalidation.Resolver,
) *HealthClient {
	if resolver == nil {
		resolver = defaultHealthResolver{}
	}

	transport := &http.Transport{
		Proxy: nil,

		TLSClientConfig: &tls.Config{
			MinVersion: tls.VersionTLS12,
		},

		DialContext: func(
			ctx context.Context,
			network string,
			address string,
		) (net.Conn, error) {
			return dialValidated(
				ctx,
				network,
				address,
				resolver,
			)
		},

		TLSHandshakeTimeout:   5 * time.Second,
		ResponseHeaderTimeout: 5 * time.Second,
		IdleConnTimeout:       30 * time.Second,
	}

	return &HealthClient{
		client: &http.Client{
			Transport: transport,
			Timeout:   healthRequestTimeout,
			CheckRedirect: func(
				req *http.Request,
				via []*http.Request,
			) error {
				return http.ErrUseLastResponse
			},
		},
		resolver: resolver,
		validate: true,
	}
}

func NewHealthClientWithHTTPClient(
	client *http.Client,
) *HealthClient {
	return &HealthClient{
		client:   client,
		validate: false,
	}
}

func (c *HealthClient) Get(
	ctx context.Context,
	rawURL string,
) (*http.Response, error) {
	if c.validate {
		if err := urlvalidation.ValidateContextWithResolver(
			ctx,
			rawURL,
			c.resolver,
		); err != nil {
			return nil, fmt.Errorf(
				"validate destination: %w",
				err,
			)
		}
	}

	request, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		rawURL,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"create health request: %w",
			err,
		)
	}

	request.Header.Set(
		"User-Agent",
		"LinkForge-Health/1.0",
	)

	return c.client.Do(request)
}

type defaultHealthResolver struct{}

func (defaultHealthResolver) LookupIP(
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

func dialValidated(
	ctx context.Context,
	network string,
	address string,
	resolver urlvalidation.Resolver,
) (net.Conn, error) {
	host, port, err := net.SplitHostPort(address)
	if err != nil {
		return nil, fmt.Errorf(
			"split destination address: %w",
			err,
		)
	}

	ips, err := resolver.LookupIP(
		ctx,
		"ip",
		host,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"resolve destination: %w",
			err,
		)
	}

	for _, ip := range ips {
		if urlvalidation.IsPrivateOrReservedIP(ip) {
			continue
		}

		dialer := net.Dialer{
			Timeout: 5 * time.Second,
		}

		conn, err := dialer.DialContext(
			ctx,
			network,
			net.JoinHostPort(
				ip.String(),
				port,
			),
		)
		if err == nil {
			return conn, nil
		}
	}

	return nil, fmt.Errorf(
		"no safe destination address available",
	)
}
