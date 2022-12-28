//go:build wireinject
// +build wireinject

package main

import (
	"context"
	"github.com/bmizerany/pat"
	"github.com/go-playground/validator/v10"
	"github.com/google/wire"

	"go_rest/internal/config"
	"go_rest/internal/rest"
	m "go_rest/internal/services/migrations"
	task_service "go_rest/internal/services/tasks"
	user_service "go_rest/internal/services/users"
	"go_rest/internal/sqlitedb"
	"net/http"
)

func initializeServer(ctx context.Context) (*http.Server, error) {
	wire.Build(
		validator.New,
		pat.New,
		config.NewConfig,
		sqlitedb.NewConnDB,
		task_service.TaskSet,
		user_service.UserSet,
		rest.NewTaskServer,
		rest.NewHttpServer,
	)
	return &http.Server{}, nil
}

func initializeMigrationService(ctx context.Context) (*m.MigrationService, error) {
	wire.Build(
		config.NewConfig,
		sqlitedb.NewConnDB,
		m.NewMigrationService,
	)
	return &m.MigrationService{}, nil
}
