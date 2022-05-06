package main

import (
	"go_rest/internal/rest"
	"log"
)

func main() {
	log.Printf("Сервер запущен")
	rest.NewTaskServer()
}
