package main

import (
	"context"

	"github.com/Clever/breakdown/gen-go/models"
	"github.com/Clever/breakdown/gen-go/server"
)

// MyController implements server.Controller
type MyController struct {
	launchConfig LaunchConfig
}

var _ server.Controller = MyController{}

// HealthCheck handles GET requests to /_health
func (mc MyController) HealthCheck(ctx context.Context) error {
	return nil
}

func (mc MyController) PostCustom(ctx context.Context, i *models.CustomData) error {
	return nil
}

func (mc MyController) PostUpload(ctx context.Context, i *models.RepoPackageFiles) error {
	return nil
}
