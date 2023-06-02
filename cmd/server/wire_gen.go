// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/omalloc/contrib/kratos/health"
	"github.com/omalloc/contrib/kratos/registry"
	"github.com/omalloc/kratos-layout/internal/biz"
	"github.com/omalloc/kratos-layout/internal/conf"
	"github.com/omalloc/kratos-layout/internal/data"
	"github.com/omalloc/kratos-layout/internal/server"
	"github.com/omalloc/kratos-layout/internal/service"
)

// Injectors from wire.go:

// wireApp init kratos application.
func wireApp(bootstrap *conf.Bootstrap, confServer *conf.Server, confData *conf.Data, logger log.Logger) (*kratos.App, func(), error) {
	protobufRegistry := server.NewRegistryConfig(bootstrap)
	client, cleanup, err := registry.NewEtcd(protobufRegistry)
	if err != nil {
		return nil, nil, err
	}
	registrar := registry.NewRegistrar(client)
	dataData, cleanup2, err := data.NewData(confData, logger)
	if err != nil {
		cleanup()
		return nil, nil, err
	}
	greeterRepo := data.NewGreeterRepo(dataData, logger)
	greeterUsecase := biz.NewGreeterUsecase(greeterRepo, logger)
	greeterService := service.NewGreeterService(greeterUsecase)
	grpcServer := server.NewGRPCServer(confServer, greeterService, logger)
	httpServer := server.NewHTTPServer(confServer, greeterService, logger)
	v := server.NewChecker(dataData)
	healthServer := health.NewServer(logger, httpServer, v)
	app := newApp(logger, registrar, grpcServer, httpServer, healthServer)
	return app, func() {
		cleanup2()
		cleanup()
	}, nil
}
