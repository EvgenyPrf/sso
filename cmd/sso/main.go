package main

import (
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	app "sso/internal/app"
	"sso/internal/config"
	"syscall"
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

	//TODO: запустить gRPC-сервер приложения
	go application.GRPCServer.MustRun()

	//Graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	//ждем, пока что-то не запишется в канал
	sign := <-stop

	log.Info("stopping application", slog.String("signal", sign.String()))

	application.GRPCServer.Stop()

	log.Info("Application stopped")

	//TODO: аналогичный graceful shutdown нужен для клиента БД и для любых воркеров (типо крона)

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
