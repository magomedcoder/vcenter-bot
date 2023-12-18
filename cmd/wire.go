//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	"vcenter-bot/internal"
)

type AppProvider struct {
	Bot *internal.Bot
}

var providerSet = wire.NewSet(
	wire.Struct(new(AppProvider), "*"),
	internal.NewBotAPI,
	internal.NewDatabase,
	internal.NewBotHandler,
	internal.NewVmwareApiCallHandler,
)

func Initialize(conf *internal.Config) *AppProvider {
	panic(wire.Build(providerSet))
}
