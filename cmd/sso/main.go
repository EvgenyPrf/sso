package main

import (
	"fmt"
	"log/slog"
	"os"
	app "sso/internal/app"
	"sso/internal/config"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	//инициализировать объект конфига
	cfg := config.MustLoad()

	fmt.Println(cfg)

	//инициализировать логгер
	log := setupLogger(cfg.Env)

	//иметь ввиду, чтобы не выводить сесурные данные
	log.Info("starting application", slog.Any("cfg", cfg))

	//TODO: инициализировать приложение

	application := app.New(log, cfg.GRPC.Port, cfg.StoragePath, cfg.TokenTTL)

	application.GRPCServer.MustRun()
	//TODO: запустить gRPC-сервер приложения
}

func setupLogger(env string) *slog.Logger {

	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(
			//для создания принимает какой-то handler
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envDev:
		log = slog.New(
			//для создания принимает какой-то handler
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			//для создания принимает какой-то handler, для прода LevelInfo
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}
