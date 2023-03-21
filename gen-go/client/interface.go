package client

import (
	"context"

	"github.com/Clever/breakdown/gen-go/models"
)

//go:generate mockgen -source=$GOFILE -destination=mock_client.go -package client --build_flags=--mod=mod -imports=models=github.com/Clever/breakdown/gen-go/models

// Client defines the methods available to clients of the breakdown service.
type Client interface {

	// HealthCheck makes a GET request to /_health
	// Checks if the service is healthy
	// 200: nil
	// 400: *models.BadRequest
	// 500: *models.InternalError
	// default: client side HTTP errors, for example: context.DeadlineExceeded.
	HealthCheck(ctx context.Context) error

	// GetThings makes a GET request to /v2/things
	//
	// 200: []models.Thing
	// 400: *models.BadRequest
	// 500: *models.InternalError
	// default: client side HTTP errors, for example: context.DeadlineExceeded.
	GetThings(ctx context.Context) ([]models.Thing, error)

	// DeleteThing makes a DELETE request to /v2/things/{id}
	//
	// 200: nil
	// 400: *models.BadRequest
	// 404: *models.NotFound
	// 500: *models.InternalError
	// default: client side HTTP errors, for example: context.DeadlineExceeded.
	DeleteThing(ctx context.Context, id string) error

	// GetThing makes a GET request to /v2/things/{id}
	//
	// 200: *models.Thing
	// 400: *models.BadRequest
	// 404: *models.NotFound
	// 500: *models.InternalError
	// default: client side HTTP errors, for example: context.DeadlineExceeded.
	GetThing(ctx context.Context, id string) (*models.Thing, error)

	// CreateOrUpdateThing makes a PUT request to /v2/things/{id}
	//
	// 200: *models.Thing
	// 400: *models.BadRequest
	// 500: *models.InternalError
	// default: client side HTTP errors, for example: context.DeadlineExceeded.
	CreateOrUpdateThing(ctx context.Context, i *models.CreateOrUpdateThingInput) (*models.Thing, error)
}
