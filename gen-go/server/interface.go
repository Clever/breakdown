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

	// PostCustom handles PUT requests to /v1/custom
	// upload or replace custom data for a given repo and commit SHA
	// 200: nil
	// 400: *models.BadRequest
	// 500: *models.InternalError
	// default: client side HTTP errors, for example: context.DeadlineExceeded.
	PostCustom(ctx context.Context, i *models.CustomData) error

	// PostDeploy handles POST requests to /v1/deploy
	// report a number of deploys
	// 200: nil
	// 400: *models.BadRequest
	// 500: *models.InternalError
	// default: client side HTTP errors, for example: context.DeadlineExceeded.
	PostDeploy(ctx context.Context, i *models.Deploys) error

	// PostUpload handles POST requests to /v1/upload
	// upload a package-type file, generated by breakdown-cli
	// 200: nil
	// 400: *models.BadRequest
	// 500: *models.InternalError
	// default: client side HTTP errors, for example: context.DeadlineExceeded.
	PostUpload(ctx context.Context, i *models.RepoCommit) error
}