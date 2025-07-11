package app

import (
	"fmt"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	ssogrpc "github.com/m1al04949/url-shortener/internal/clients/sso/grpc"
	"github.com/m1al04949/url-shortener/internal/config"
	"github.com/m1al04949/url-shortener/internal/http-server/handlers/delete"
	"github.com/m1al04949/url-shortener/internal/http-server/handlers/url/save"
	mwLogger "github.com/m1al04949/url-shortener/internal/http-server/middleware/logger"
	"github.com/m1al04949/url-shortener/internal/lib/logger/logslog"
	"github.com/m1al04949/url-shortener/internal/pkg/setlog"
	"github.com/m1al04949/url-shortener/internal/storage/sqlite"
	httpSwagger "github.com/swaggo/http-swagger"
	"golang.org/x/exp/slog"
	"golang.org/x/net/context"
)

func RunServer() error {
	// Init Config
	cfg := config.MustLoad()
	//Init Logger
	log := setlog.SetupLogger(cfg.Env)

	log.Info("starting url-shortener", slog.String("env", cfg.Env))
	log.Debug("debug messages are enabled")
	// Init Storage
	storage, err := sqlite.New(cfg.StoragePath)
	if err != nil {
		log.Error("failed to init storage", logslog.Err(err))
		return err
	}
	// Init Clients
	ssoClient, err := ssogrpc.New(
		context.Background(),
		log,
		cfg.Clients.SSO.Address,
		cfg.Clients.SSO.Timeout,
		cfg.Clients.SSO.RetriesCount,
	)
	if err != nil {
		log.Error("failed to init sso client", logslog.Err(err))
		os.Exit(1)
	}
	// Check from sso client (gRPC)
	ssoClient.IsAdmin(context.Background(), 1)

	// Init Router
	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(mwLogger.New(log))
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	router.Route("/url", func(r chi.Router) {
		r.Use(middleware.BasicAuth("url-shortener", map[string]string{
			cfg.HTTPServer.User: cfg.HTTPServer.Password,
		}))

		r.Post("/", save.New(log, storage))
		r.Delete("/{alias}", delete.New(log, storage))
	})
	// Swagger UI
	router.Get("/swagger/*", httpSwagger.WrapHandler)

	log.Info("starting server", slog.String("address", cfg.Address))
	// Init Server
	srv := &http.Server{
		Addr:         cfg.Address,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Error("failed to start server")
		return err
	}

	log.Error("server stopped")

	return fmt.Errorf("server stopped")
}
