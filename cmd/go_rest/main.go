package main

import (
	"context"
	"fmt"
	"go_rest/internal/config"
	"go_rest/internal/logger"
	"go_rest/internal/rest"
	"go_rest/internal/taskstore/sqlitestore"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"golang.org/x/sync/errgroup"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)

		<-c
		cancel()
	}()

	ts, err := sqlitestore.New("test.db")
	if err != nil {
		log.Fatal(err)
	}

	url := fmt.Sprintf("localhost:%v", config.Config.Port)
	s := rest.NewTaskServer(ts)

	httpServer := &http.Server{
		Addr:    url,
		Handler: s.Router,
	}

	g, gCtx := errgroup.WithContext(ctx)
	g.Go(func() error {
		logger.Info("Сервер запущен")
		err = httpServer.ListenAndServe()
		if err != nil {
			return err
		}
		return nil
	})
	g.Go(func() error {
		<-gCtx.Done()
		ts.Close()
		return httpServer.Shutdown(context.Background())
	})

	if err := g.Wait(); err != nil {
		logger.Info(fmt.Sprintf("Сервер остановлен : %s \n", err))

	}

}
