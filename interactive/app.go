package main

import (
	"github.com/wx-up/go-book/pkg/grpcx"
	"github.com/wx-up/go-book/pkg/saramax"
)

type App struct {
	consumers []saramax.Consumer
	server    *grpcx.Server
}
