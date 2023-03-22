package main

import (
	"context"
	"flag"
	"log"
	"os"
	"time"

	"github.com/Clever/breakdown/gen-go/server"
	"github.com/Clever/breakdown/gen-go/servertracing"
	"github.com/Clever/kayvee-go/v7/logger"
	"github.com/Clever/kayvee-go/v7/middleware"
	"github.com/Clever/wag/swagger"
	trace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
)

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
