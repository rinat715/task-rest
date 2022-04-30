package main

import (
	"go_rest/internal/logger"
	"go_rest/internal/rest"
)

func main() {
	logger.Info("Сервер запущен")
	rest.NewTaskServer()
}
