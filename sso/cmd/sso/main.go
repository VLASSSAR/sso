package main

import (
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"sso/internal/app"
	"sso/internal/config"
	"sso/internal/lib/logger/handlers/slogpretty"
	"sso/internal/lib/logger/sl"
	"syscall"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {

	// TODO: инициализировать объект конфига
	cfg := config.MustLoad()

	fmt.Println(cfg)

	// TODO: инициализировать logger
	log := setupLogger(cfg.Env)

	log.Info("starting application", slog.Any("config", cfg))
	//slog.String("env", cfg.Env),
	//slog.Any("cfg", cfg),
	//slog.Int("port", cfg.GRPC.Port),
	//)
	//log.Debug("debug message")
	//
	//log.Error("error message")
	//
	//log.Warn("warn message")

	var err error

	if err != nil {
		//log.Error("error message", slog.String("error", err.Error()))
		log.Error("error message", sl.Err(err))
	}

	// TODO: инициализировать приложение (app)
	application := app.New(log, cfg.GRPC.Port, cfg.StoragePath, cfg.TokenTTL)

	go application.GRPCSrv.MustRun()

	// TODO: запустить gRPC-сервер приложение

	// TODO: Graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	sign := <-stop

	log.Info("stopping application", slog.String("signal", sign.String()))

	application.GRPCSrv.Stop()

	log.Info("application stopped")
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = setupPrettySlog()
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

func setupPrettySlog() *slog.Logger {
	opts := slogpretty.PrettyHandlerOptions{
		SlogOpts: &slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	}

	handler := opts.NewPrettyHandler(os.Stdout)

	return slog.New(handler)
}
