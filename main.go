package main

import (
	"context"
	"flag"
	"log"
	"os"
	"time"

	"github.com/Clever/breakdown/gen-go/models"
	"github.com/Clever/breakdown/gen-go/server"
	"github.com/Clever/breakdown/gen-go/servertracing"
	"github.com/Clever/kayvee-go/v7/logger"
	"github.com/Clever/kayvee-go/v7/middleware"
	"github.com/Clever/wag/swagger"
	trace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
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

// GetThings handles GET requests to /v2/things
func (mc MyController) GetThings(ctx context.Context) ([]models.Thing, error) {
	return nil, models.InternalError{Message: "TODO, this is just an example"}
}

// DeleteThing handles DELETE requests to /v2/things/{id}
func (mc MyController) DeleteThing(ctx context.Context, id string) error {
	return models.InternalError{Message: "TODO, this is just an example"}
}

// GetThing handles GET requests to /v2/things/{id}
func (mc MyController) GetThing(ctx context.Context, id string) (*models.Thing, error) {
	return nil, models.NotFound{Message: "TODO, this is just an example"}
}

// CreateOrUpdateThing handles PUT requests to /v2/things/{id}
func (mc MyController) CreateOrUpdateThing(ctx context.Context, i *models.CreateOrUpdateThingInput) (*models.Thing, error) {
	return nil, models.InternalError{Message: "TODO, this is just an example"}
}

func main() {
	addr := flag.String("addr", ":8080", "Address to listen at")
	flag.Parse()

	swagger.InitCustomFormats()

	// Initialize globals for tracing
	var exporter trace.SpanExporter = tracetest.NewNoopExporter()
	if os.Getenv("_TRACING_ENABLED") == "true" {
		exp, prov, err := servertracing.SetupGlobalTraceProviderAndExporter(context.Background())
		if err != nil {
			log.Fatalf("failed to setup tracing: %v", err)
		}
		exporter = exp
		// Ensure traces are finalized when exiting
		defer exp.Shutdown(context.Background())
		defer prov.Shutdown(context.Background())
	}

	middleware.EnableRollups(context.Background(), logger.NewConcreteLogger("breakdown"), 20*time.Second)

	myController := MyController{
		launchConfig: InitLaunchConfig(&exporter),
	}
	s := server.New(myController, *addr)

	if err := s.Serve(); err != nil {
		log.Fatal(err)
	}

	log.Println("breakdown exited without error")
}
