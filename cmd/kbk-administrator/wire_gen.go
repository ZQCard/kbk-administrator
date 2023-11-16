// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"github.com/ZQCard/kbk-administrator/internal/biz"
	"github.com/ZQCard/kbk-administrator/internal/conf"
	"github.com/ZQCard/kbk-administrator/internal/data"
	"github.com/ZQCard/kbk-administrator/internal/server"
	"github.com/ZQCard/kbk-administrator/internal/service"
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
)

import (
	_ "go.uber.org/automaxprocs"
)

// Injectors from wire.go:

// wireApp init kratos application.
func wireApp(env *conf.Env, confServer *conf.Server, confData *conf.Data, bootstrap *conf.Bootstrap, logger log.Logger) (*kratos.App, func(), error) {
	db := data.NewMysqlCmd(bootstrap, logger)
	client := data.NewRedisClient(confData)
	dataData, cleanup, err := data.NewData(bootstrap, db, client, logger)
	if err != nil {
		return nil, nil, err
	}
	administratorRepo := data.NewAdministratorRepo(dataData, logger)
	administratorUsecase := biz.NewAdministratorUsecase(administratorRepo, logger)
	administratorService := service.NewAdministratorService(administratorUsecase, logger)
	grpcServer := server.NewGRPCServer(confServer, administratorService, logger)
	httpServer := server.NewHTTPServer(confServer, administratorService, logger)
	app := newApp(logger, grpcServer, httpServer)
	return app, func() {
		cleanup()
	}, nil
}
