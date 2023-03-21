package server

import (
	"context"

	"github.com/Clever/breakdown/gen-go/models"
)

//go:generate mockgen -source=$GOFILE -destination=mock_controller.go -package server --build_flags=--mod=mod -imports=models=github.com/Clever/breakdown/gen-go/models

// Controller defines the interface for the breakdown service.
type Controller interface {

	// HealthCheck handles GET requests to /_health
	// Checks if the service is healthy
	// 200: nil
	// 400: *models.BadRequest
	// 500: *models.InternalError
	// default: client side HTTP errors, for example: context.DeadlineExceeded.
	HealthCheck(ctx context.Context) error

	// GetThings handles GET requests to /v2/things
	//
	// 200: []models.Thing
	// 400: *models.BadRequest
	// 500: *models.InternalError
	// default: client side HTTP errors, for example: context.DeadlineExceeded.
	GetThings(ctx context.Context) ([]models.Thing, error)

	// DeleteThing handles DELETE requests to /v2/things/{id}
	//
	// 200: nil
	// 400: *models.BadRequest
	// 404: *models.NotFound
	// 500: *models.InternalError
	// default: client side HTTP errors, for example: context.DeadlineExceeded.
	DeleteThing(ctx context.Context, id string) error

	// GetThing handles GET requests to /v2/things/{id}
	//
	// 200: *models.Thing
	// 400: *models.BadRequest
	// 404: *models.NotFound
	// 500: *models.InternalError
	// default: client side HTTP errors, for example: context.DeadlineExceeded.
	GetThing(ctx context.Context, id string) (*models.Thing, error)

	// CreateOrUpdateThing handles PUT requests to /v2/things/{id}
	//
	// 200: *models.Thing
	// 400: *models.BadRequest
	// 500: *models.InternalError
	// default: client side HTTP errors, for example: context.DeadlineExceeded.
	CreateOrUpdateThing(ctx context.Context, i *models.CreateOrUpdateThingInput) (*models.Thing, error)
}
