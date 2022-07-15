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
		logger.Info("миграции")
		m, err := initializeMigrationService(ctx)
		if err != nil {
			panic(err)
		}
		err = m.Make(*migratePath)
		if err != nil {
			panic(err)
		}
	} else {
		logger.Info("старт сервера")
		s, err := initializeServer(ctx)
		if err != nil {
			panic(err)
		}
		logger.Info("Сервер запущен")
		go func() {
			<-ctx.Done()
			err := s.Shutdown(context.Background())
			logger.Info(fmt.Sprintf("Сервер остановлен : %s", err))
		}()

		err = s.ListenAndServe()
		if err != nil {
			logger.Info(fmt.Sprintf("Сервер не запущен : %s", err))
		}

	}

}
