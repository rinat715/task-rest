package main

import (
	"context"
	"fmt"
	"go_rest/internal/logger"
	"os"
	"os/signal"
	"syscall"

	"golang.org/x/sync/errgroup"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	s, err := initializeServer(ctx)
	if err != nil {
		panic(err)
	}

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)

		<-c
		cancel()
	}()

	g, gCtx := errgroup.WithContext(ctx)
	g.Go(func() error {
		logger.Info("Сервер запущен")
		err = s.ListenAndServe()
		if err != nil {
			return err
		}
		return nil
	})
	g.Go(func() error {
		<-gCtx.Done()
		return s.Shutdown(context.Background())
	})

	if err := g.Wait(); err != nil {
		logger.Info(fmt.Sprintf("Сервер остановлен : %s \n", err))

	}
}
