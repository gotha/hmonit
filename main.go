package main

import (
	"context"
	"os"
	"time"

	"github.com/Financial-Times/gourmet/app"
	"github.com/Financial-Times/gourmet/log"
	ghttp "github.com/Financial-Times/gourmet/transport/http"
	"github.com/go-chi/chi"
)

func main() {
	conf := getConfig()
	logger := log.NewStructuredLogger(
		conf.GetLogLevel(),
		log.WithServiceName(conf.SystemCode),
	)

	hs := NewHealthStore()
	hc := NewHealthChecker()
	m := &HealthMonitor{
		checker:     hc,
		healthStore: hs,
		services:    conf.Services,
		logger:      logger,
	}
	healthcheckerLifecycle := NewHealthCheckLifescycle(m, conf.RefreshInterval)

	rlm := NewRequestLoggingMiddleware(logger)
	hh := NewHealthHandler(hs, logger)

	r := chi.NewRouter()
	r.Route("/", func(r chi.Router) {
		r.Get("/", rlm.Log(hh.Status))
	})

	srv := ghttp.New(r,
		ghttp.WithCustomAppPort(conf.Server.Port),
		ghttp.WithCustomHTTPServerReadTimeout(conf.Server.ReadTimeout),
		ghttp.WithCustomHTTPServerWriteTimeout(conf.Server.WriteTimeout),
		ghttp.WithCustomHTTPServerIdleTimeout(conf.Server.IdleTimeout),
	)

	httpLifecycle := ghttp.NewServerLifecycle(srv)

	a := app.New([]app.Lifecycle{
		httpLifecycle,
		healthcheckerLifecycle,
	})

	logger.Info("Service starting",
		log.WithField("appPort", conf.Server.Port),
	)
	startCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	if err := a.Run(startCtx); err != nil {
		logger.Error("error running application", log.WithError(err))
		os.Exit(1)
	}
}
