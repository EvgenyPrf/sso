package grpc

import (
	"log/slog"
	grpcapp "sso/internal/app/grpc"
	auth "sso/internal/services/auth"
	storage "sso/internal/storage/sqlite"
	"time"
)

type App struct {
	GRPCServer *grpcapp.App
}

func New(log *slog.Logger, grpcPort int, storagePath string, tokenTtl time.Duration) *App {

	//инициализировать хранилище
	strg, err := storage.New(storagePath)

	if err != nil {
		panic(err)
	}

	//инициализировать сервисный слой auth сервиса
	authService := auth.New(log, strg, strg, strg, tokenTtl)

	grpcApp := grpcapp.New(log, authService, grpcPort)

	return &App{
		GRPCServer: grpcApp,
	}
}
