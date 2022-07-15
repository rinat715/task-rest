package main

import (
	"context"
	"flag"
	"fmt"
	"go_rest/internal/logger"
	"os"
	"os/signal"
	"syscall"
)

var (
	migratePath *string
)

func init() {
	migratePath = flag.String("migratePath", "", "Путь к каталогу с миграциями")
}

func runServer(ctx context.Context) {
	s, err := initializeServer(ctx)
	if err != nil {
		panic(err)
	}
	go func() {
		logger.Info("Сервер запущен")
		err := s.ListenAndServe()
		if err != nil {
			panic(err)
		}
	}()
	go func() {
		<-ctx.Done()
		err := s.Shutdown(context.Background())
		logger.Info(fmt.Sprintf("Сервер остановлен : %s \n", err))
	}()
}

func main() {
	flag.Parse()

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)

		<-c
		cancel()
	}()

	if *migratePath != "" {
		m, err := initializeMigrationService(ctx)
		if err != nil {
			panic(err)
		}
		m.Make(*migratePath)
	} else {
		runServer(ctx)
	}

}
