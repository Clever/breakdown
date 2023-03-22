package main

import (
	"context"
	"flag"
	"log"
	"os"
	"time"

	"github.com/Clever/breakdown/db"
	"github.com/Clever/breakdown/gen-go/server"
	"github.com/Clever/breakdown/gen-go/servertracing"
	"github.com/Clever/kayvee-go/v7/logger"
	"github.com/Clever/kayvee-go/v7/middleware"
	"github.com/Clever/wag/swagger"
	"github.com/jackc/pgx/v4"
	trace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
)

func newDBFromLaunch(config LaunchConfig) (*pgx.Conn, error) {
	return db.FromConfig(db.Config{
		User:         config.Env.PostgresUsername,
		Password:     config.Env.PostgresPassword,
		Host:         config.Env.PostgresHost,
		DatabaseName: config.Env.PostgresDb,
		Port:         "5432",
	})
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

	launchConfig := InitLaunchConfig(&exporter)

	db, err := newDBFromLaunch(launchConfig)
	if err != nil {
		log.Fatal(err)
	}

	myController := MyController{
		launchConfig: launchConfig,
		db:           db,
		l:            logger.NewConcreteLogger("breakdown"),
	}
	s := server.New(myController, *addr)

	if err := s.Serve(); err != nil {
		log.Fatal(err)
	}

	log.Println("breakdown exited without error")
}
