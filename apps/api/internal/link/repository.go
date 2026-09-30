package link

import (
	"context"
	"errors"
)

var (
	ErrLinkNotFound   = errors.New("link not found")
	ErrShortCodeTaken = errors.New("short code already exists")
)

type Repository interface {
	Create(
		ctx context.Context,
		link *Link,
	) error

	GetByShortCode(
		ctx context.Context,
		shortCode string,
	) (*Link, error)

	List(
		ctx context.Context,
	) ([]*Link, error)

	UpdateStatus(
		ctx context.Context,
		shortCode string,
		status Status,
	) error

	Delete(
		ctx context.Context,
		shortCode string,
	) error
}
