package link

import (
	"context"
	"errors"
	"strings"
	"time"

	urlvalidation "github.com/sm-joe/linkforge/packages/url-validation"
)

var (
	ErrDestinationRequired = errors.New("destination URL is required")
	ErrInvalidAlias        = errors.New("invalid alias")
	ErrAliasTooLong        = errors.New("alias is too long")

	ErrDestinationPrivateAddress = errors.New("destination resolves to a private or reserved address")

	ErrDestinationDNSResolution = errors.New("destination DNS resolution failed")
)

const (
	maxAliasLength = 64
)

type CreateRequest struct {
	Destination string
	Alias       string
	ExpiresAt   *time.Time
}

type Service struct {
	repository Repository
	resolver   urlvalidation.Resolver
}

func NewService(
	repository Repository,
	resolvers ...urlvalidation.Resolver,
) *Service {
	var resolver urlvalidation.Resolver

	if len(resolvers) > 0 {
		resolver = resolvers[0]
	}

	return &Service{
		repository: repository,
		resolver:   resolver,
	}
}

func (s *Service) ValidateDestination(
	destination string,
) error {
	return s.ValidateDestinationContext(
		context.Background(),
		destination,
	)
}

func (s *Service) ValidateDestinationContext(
	ctx context.Context,
	destination string,
) error {
	destination = strings.TrimSpace(destination)

	if destination == "" {
		return ErrDestinationRequired
	}

	var err error

	if s.resolver != nil {
		err = urlvalidation.ValidateContextWithResolver(
			ctx,
			destination,
			s.resolver,
		)
	} else {
		err = urlvalidation.ValidateContext(
			ctx,
			destination,
		)
	}

	switch {
	case errors.Is(
		err,
		urlvalidation.ErrPrivateAddress,
	):
		return ErrDestinationPrivateAddress

	case errors.Is(
		err,
		urlvalidation.ErrDNSResolution,
	):
		return ErrDestinationDNSResolution

	default:
		return err
	}
}

func (s *Service) CreateLink(
	ctx context.Context,
	request CreateRequest,
) (*Link, error) {
	destination := strings.TrimSpace(
		request.Destination,
	)

	if err := s.ValidateDestinationContext(
		ctx,
		destination,
	); err != nil {
		return nil, err
	}

	alias := strings.TrimSpace(request.Alias)

	if len(alias) > maxAliasLength {
		return nil, ErrAliasTooLong
	}

	if alias != "" && !isValidAlias(alias) {
		return nil, ErrInvalidAlias
	}

	shortCode := alias

	if shortCode == "" {
		shortCode = generateShortCode()
	}

	now := time.Now().UTC()

	link := &Link{
		ID:          generateID(),
		ShortCode:   shortCode,
		Alias:       alias,
		Destination: destination,
		Status:      StatusActive,
		CreatedAt:   now,
		ExpiresAt:   request.ExpiresAt,
	}

	if err := s.repository.Create(ctx, link); err != nil {
		return nil, err
	}

	return link, nil
}

func isValidAlias(alias string) bool {
	for _, r := range alias {
		if (r >= 'a' && r <= 'z') ||
			(r >= 'A' && r <= 'Z') ||
			(r >= '0' && r <= '9') ||
			r == '-' ||
			r == '_' {
			continue
		}

		return false
	}

	return true
}

func (s *Service) GetByShortCode(
	ctx context.Context,
	shortCode string,
) (*Link, error) {
	shortCode = strings.TrimSpace(shortCode)

	if shortCode == "" {
		return nil, ErrLinkNotFound
	}

	return s.repository.GetByShortCode(ctx, shortCode)
}

func (s *Service) DisableLink(
	ctx context.Context,
	shortCode string,
) error {
	shortCode = strings.TrimSpace(shortCode)

	if shortCode == "" {
		return ErrLinkNotFound
	}

	return s.repository.UpdateStatus(
		ctx,
		shortCode,
		StatusDisabled,
	)
}

func (s *Service) EnableLink(
	ctx context.Context,
	shortCode string,
) error {
	shortCode = strings.TrimSpace(shortCode)

	if shortCode == "" {
		return ErrLinkNotFound
	}

	return s.repository.UpdateStatus(
		ctx,
		shortCode,
		StatusActive,
	)
}

func (s *Service) DeleteLink(
	ctx context.Context,
	shortCode string,
) error {
	shortCode = strings.TrimSpace(shortCode)

	if shortCode == "" {
		return ErrLinkNotFound
	}

	return s.repository.Delete(
		ctx,
		shortCode,
	)
}

func (s *Service) ListLinks(
	ctx context.Context,
) ([]*Link, error) {
	return s.repository.List(ctx)
}
