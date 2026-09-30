package link

import (
	"context"
	"testing"
	"time"
)

type fakeRepository struct {
	created *Link
}

func (f *fakeRepository) Create(
	_ context.Context,
	link *Link,
) error {
	f.created = link
	return nil
}

func (f *fakeRepository) GetByShortCode(
	_ context.Context,
	_ string,
) (*Link, error) {
	return nil, ErrLinkNotFound
}

func TestCreateLink(t *testing.T) {
	repository := &fakeRepository{}
	service := NewService(repository)

	link, err := service.CreateLink(
		context.Background(),
		CreateRequest{
			Destination: "https://example.com/test",
		},
	)

	if err != nil {
		t.Fatalf("CreateLink() error = %v", err)
	}

	if link == nil {
		t.Fatal("CreateLink() returned nil link")
	}

	if link.ShortCode == "" {
		t.Fatal("expected short code")
	}

	if len(link.ShortCode) != shortCodeLength {
		t.Fatalf(
			"short code length = %d, want %d",
			len(link.ShortCode),
			shortCodeLength,
		)
	}

	if link.Status != StatusActive {
		t.Fatalf(
			"status = %q, want %q",
			link.Status,
			StatusActive,
		)
	}
}

func TestCreateLinkWithAlias(t *testing.T) {
	repository := &fakeRepository{}
	service := NewService(repository)

	link, err := service.CreateLink(
		context.Background(),
		CreateRequest{
			Destination: "https://example.com",
			Alias:       "my-project",
		},
	)

	if err != nil {
		t.Fatalf("CreateLink() error = %v", err)
	}

	if link.ShortCode != "my-project" {
		t.Fatalf(
			"short code = %q, want %q",
			link.ShortCode,
			"my-project",
		)
	}
}

func TestCreateLinkRejectsInvalidAlias(t *testing.T) {
	repository := &fakeRepository{}
	service := NewService(repository)

	_, err := service.CreateLink(
		context.Background(),
		CreateRequest{
			Destination: "https://example.com",
			Alias:       "bad alias!",
		},
	)

	if err != ErrInvalidAlias {
		t.Fatalf(
			"error = %v, want %v",
			err,
			ErrInvalidAlias,
		)
	}
}

func TestCreateLinkSupportsExpiration(t *testing.T) {
	repository := &fakeRepository{}
	service := NewService(repository)

	expiresAt := time.Now().UTC().Add(24 * time.Hour)

	link, err := service.CreateLink(
		context.Background(),
		CreateRequest{
			Destination: "https://example.com",
			ExpiresAt:   &expiresAt,
		},
	)

	if err != nil {
		t.Fatalf("CreateLink() error = %v", err)
	}

	if link.ExpiresAt == nil {
		t.Fatal("expected expiration time")
	}
}

func (f *fakeRepository) UpdateStatus(
	_ context.Context,
	_ string,
	_ Status,
) error {
	return nil
}

func (f *fakeRepository) Delete(
	_ context.Context,
	_ string,
) error {
	return nil
}

func TestLinkExpiration(t *testing.T) {
	repository := &fakeRepository{}
	service := NewService(repository)

	expired := time.Now().UTC().Add(-time.Hour)

	link, err := service.CreateLink(
		context.Background(),
		CreateRequest{
			Destination: "https://example.com",
			ExpiresAt:   &expired,
		},
	)

	if err != nil {
		t.Fatalf("CreateLink() error = %v", err)
	}

	if !link.IsExpired(time.Now().UTC()) {
		t.Fatal("expected link to be expired")
	}

	if link.IsAvailable(time.Now().UTC()) {
		t.Fatal("expected expired link to be unavailable")
	}
}

func TestDisabledLinkIsUnavailable(t *testing.T) {
	link := &Link{
		Status: StatusDisabled,
	}

	if link.IsAvailable(time.Now().UTC()) {
		t.Fatal("expected disabled link to be unavailable")
	}
}

func TestActiveLinkIsAvailable(t *testing.T) {
	link := &Link{
		Status: StatusActive,
	}

	if !link.IsAvailable(time.Now().UTC()) {
		t.Fatal("expected active link to be available")
	}
}
