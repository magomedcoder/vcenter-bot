// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"github.com/google/wire"
	"vcenter-bot/internal"
)

// Injectors from wire.go:

func Initialize(conf *internal.Config) *AppProvider {
	botAPI := internal.NewBotAPI(conf)
	vCenterApiCall := internal.NewVmwareApiCallHandler(conf)
	bot := internal.NewBotHandler(botAPI, vCenterApiCall)
	appProvider := &AppProvider{
		Bot: bot,
	}
	return appProvider
}

// wire.go:

type AppProvider struct {
	Bot *internal.Bot
}

var providerSet = wire.NewSet(wire.Struct(new(AppProvider), "*"), internal.NewBotAPI, internal.NewBotHandler, internal.NewVmwareApiCallHandler)
