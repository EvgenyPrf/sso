package main

import (
	"fmt"
	"sso/internal/config"
)

func main() {
	//инициализировать объект конфига
	cfg := config.MustLoad()

	fmt.Println(cfg)

	//TODO: инициализировать логгер

	//TODO: инициализировать приложение (app)

	//TODO: запустить gRPC-сервер приложения
}
