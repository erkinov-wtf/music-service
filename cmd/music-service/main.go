// cmd/music-service/main.go
package main

import (
	"context"
	"go.uber.org/fx"
	"log"
	"log/slog"
	"music-service/internal/api/handlers"
	"music-service/internal/api/routes"
	"music-service/internal/api/services"
	"music-service/internal/config"
	"music-service/internal/storage/database/repository"
	"os"
	"os/signal"
	"syscall"
)

func provideEnv(cfg *config.Config) string {
	env := config.LocalEnv
	if cfg.Env != "" {
		env = cfg.Env
	}
	return env
}

func provideDBManager(cfg *config.Config) *repository.Manager {
	ctx := context.Background()
	return repository.MustConnectDB(cfg, ctx)
}

func provideRepositories(dbManager *repository.Manager) (repository.GroupRepositoryInterface, repository.SongRepositoryInterface) {
	return dbManager.Groups, dbManager.Songs
}

// Add this function to provide a *slog.Logger
func provideLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
}

func main() {
	app := fx.New(
		fx.Provide(
			config.MustLoad,
			provideEnv,

			provideLogger,

			// Database and repositories
			provideDBManager,
			provideRepositories,

			// Services
			services.NewSongService,
			services.NewGroupService,

			// Handlers setup
			handlers.NewGroupHandler,
			handlers.NewSongHandler,

			// Router
			routes.NewRouter,
		),

		// Register hooks
		fx.Invoke(routes.RegisterRoutes),

		// Lifecycle hooks
		fx.Invoke(registerHooks),
		fx.Invoke(startHTTPServer),
	)

	startCtx, cancel := context.WithTimeout(context.Background(), config.DefaultTimeout)
	defer cancel()

	if err := app.Start(startCtx); err != nil {
		log.Fatal(err)
	}

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	// Graceful shutdown
	stopCtx, cancel := context.WithTimeout(context.Background(), config.DefaultTimeout)
	defer cancel()

	if err := app.Stop(stopCtx); err != nil {
		log.Fatal(err)
	}
}

func startHTTPServer(lc fx.Lifecycle, router *routes.Router, log *slog.Logger) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			log.Info("Starting HTTP server")
			// Start the server in a goroutine so it doesn't block
			go func() {
				if err := router.Run(); err != nil {
					log.Error("Failed to start HTTP server", "error", err)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			log.Info("Stopping HTTP server")
			// The HTTP server will stop when the application stops
			return nil
		},
	})
}

func registerHooks(lc fx.Lifecycle, dbManager *repository.Manager, cfg *config.Config, log *slog.Logger) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			log.Info("Starting music service")
			return nil
		},
		OnStop: func(ctx context.Context) error {
			log.Info("Shutting down music service")
			dbManager.Close()
			return nil
		},
	})
}
