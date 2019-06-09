package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	opentracing "github.com/opentracing/opentracing-go"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
	jcfg "github.com/uber/jaeger-client-go/config"

	"github.com/cedric-parisi/fizzbuzz-api/config"
	"github.com/cedric-parisi/fizzbuzz-api/internal/fizzbuzz/repository"
	"github.com/cedric-parisi/fizzbuzz-api/internal/fizzbuzz/service"
	"github.com/cedric-parisi/fizzbuzz-api/internal/fizzbuzz/transport"
	"github.com/cedric-parisi/fizzbuzz-api/pkg/instrumenting"
	"github.com/cedric-parisi/fizzbuzz-api/pkg/tracing"
	"github.com/cedric-parisi/fizzbuzz-api/storage"
)

const (
	appName          = "fizzbuzz-api"
	fizzbuzzComp     = "fizzbuzz"
	gracefulPeriod   = 5 * time.Second
	httpReadTimeout  = 5 * time.Second
	httpWriteTimeout = 60 * time.Second
)

func main() {
	// Init default context.
	ctx, cancel := context.WithTimeout(context.Background(), gracefulPeriod)
	defer cancel()

	// Setup configuration.
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("could not load config: %v", err)
	}

	// Listen to interruption signal from the system.
	stopChan := make(chan os.Signal)
	signal.Notify(stopChan, os.Interrupt)

	// Setup postgres storage.
	pq, err := storage.NewPostgres(cfg)
	if err != nil {
		log.Fatalf("could not connect to postgres: %v", err)
	}
	defer pq.Close()

	// Setup tracing.
	jc, err := jcfg.FromEnv()
	if err != nil {
		log.Fatalf("could not load jaeger config: %v", err)
	}
	tracer, closer, err := jc.NewTracer()
	if err != nil {
		log.Fatalf("could not create tracer: %v", err)
	}
	defer closer.Close()
	opentracing.SetGlobalTracer(tracer)

	// Setup logger.
	log := logrus.New()
	log.SetFormatter(&logrus.JSONFormatter{})
	log.Out = os.Stderr

	// Setup echo router.
	router := echo.New()
	router.Pre(middleware.AddTrailingSlash())
	router.Use(middleware.Recover())

	// Setup health endpoint.
	router.GET("/healthz/", func(etx echo.Context) error {
		return etx.NoContent(http.StatusOK)
	})

	// Expose metrics endpoint.
	router.GET("/metrics/", echo.WrapHandler(promhttp.Handler()))

	// Expose documentation endpoint
	router.Static("/swaggerui", "swaggerui/")

	// Setup fizzbuzz handler.
	transport.NewFizzbuzzHandler(
		router.Group("/v1/fizzbuzz",
			instrumenting.Metrics(fizzbuzzComp),
			tracing.Tracing(fizzbuzzComp)),
		service.NewFizzbuzzService(repository.NewPostgresFizzbuzzRepository(pq, cfg.DbTimeout)),
		log,
	)

	// Init HTTP server.
	srv := &http.Server{
		Addr:         ":" + cfg.AppPort,
		ReadTimeout:  httpReadTimeout,
		WriteTimeout: httpWriteTimeout,
		Handler:      router,
	}

	go func() {
		log.Printf("%s listening on port %s", appName, cfg.AppPort)
		if err := srv.ListenAndServe(); err != nil {
			if err != http.ErrServerClosed {
				log.Fatal(err)
			}
		}
	}()

	// Block here until stop signal received
	<-stopChan

	// Graceful shutdown
	log.Print("shutting down...")
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal(err)
	}

	log.Print("gracefully stopped")
}
