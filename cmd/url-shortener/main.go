package main

import (
	"log/slog"
	"net/http"
	"os"
	"time"

	"url-shortener/internal/config"
	"url-shortener/internal/metrics"
	"url-shortener/internal/storage/postgres"

	dropdb "url-shortener/internal/http-server/handlers/db"
	"url-shortener/internal/http-server/handlers/url/save"

	mwlogger "url-shortener/internal/http-server/middleware/logger"

	"url-shortener/internal/http-server/handlers/url/delete"
	"url-shortener/internal/http-server/handlers/url/redirect"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

// bug
// bug 2
// bug 3
// TODO: add docs
func main() {

	//cfg := config.MustLoad()

	cfg := config.Config{
		Env:              "local",
		ConnectionString: "host=db port=5432 user=postgres password=230704 dbname=postgres sslmode=disable",
		HTTPServer: config.HTTPServer{
			Address:     ":8000",
			Timeout:     time.Second * 4,
			IdleTimeout: time.Minute,
			User:        "myuser",
			Password:    "mypass",
		},
		Metrics: config.Metrics{
			PrometheusAddress: ":9001",
		},
	}

	log := setupLogger(cfg.Env)
	log.Info("starting url-shortener", slog.String("env", cfg.Env))
	log.Debug("debug messages are enabled")

	storage, err := postgres.New(cfg.ConnectionString)
	if err != nil {
		panic(err.Error())
	}

	router := chi.NewRouter()

	router.Use(middleware.RequestID) // Каждому запросу добавляется его id
	router.Use(mwlogger.New(log))    // Кастомный logger middleware
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	pMux, metrics := setupPrometheus()

	router.Route("/url", func(r chi.Router) {
		r.Use(middleware.BasicAuth("url-shortener", map[string]string{
			cfg.HTTPServer.User: cfg.HTTPServer.Password,
		}))
		r.Post("/", save.New(log, storage, metrics))
		r.Get("/drop", dropdb.New(log, storage))
		r.Delete("/{alias}", delete.New(log, storage))
	})

	router.Get("/{alias}", redirect.New(log, storage, metrics))

	log.Info("starting server", slog.String("address", cfg.Address))

	srv := &http.Server{
		Addr:         cfg.Address,
		Handler:      router,
		ReadTimeout:  cfg.Timeout,
		WriteTimeout: cfg.Timeout,
		IdleTimeout:  cfg.IdleTimeout,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Error("failed to start server")
		}
	}()

	go func() {
		if err := http.ListenAndServe(cfg.PrometheusAddress, pMux); err != nil {
			log.Error("failed to start prometheus server")
		}
	}()

	select {}
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger
	switch env {
	case envLocal:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}
	return log
}

func setupPrometheus() (*http.ServeMux, *metrics.Metrics) {
	reg := prometheus.NewRegistry()

	m := metrics.NewMetrics(reg)

	promHandler := promhttp.HandlerFor(reg, promhttp.HandlerOpts{})

	pMux := http.NewServeMux()

	pMux.Handle("/metrics", promHandler)

	return pMux, m
}
