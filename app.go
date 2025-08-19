package main

import (
	"context"
	"orcrux/shamir"
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

func (a *App) Split(secret string, shards int, shardsNeeded int, output string) string {
	out, err := shamir.Split([]byte(secret), shards, shardsNeeded, output)
	if err != nil {
		return "error: " + err.Error()
	}
	return out
}

func (a *App) Recompose(shards []string) string {
	out, err := shamir.Recompose(shards)
	if err != nil {
		return "error: " + err.Error()
	}
	return string(out)
}
