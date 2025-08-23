package main

import (
	"context"
	"encoding/json"
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

// Response represents the standard response format
type Response struct {
	Error *string     `json:"error"`
	Data  interface{} `json:"data"`
}

func (a *App) Split(secret string, shards int, shardsNeeded int, output string) string {
	out, err := shamir.Split([]byte(secret), shards, shardsNeeded, output)

	response := Response{}
	if err != nil {
		errorMsg := err.Error()
		response.Error = &errorMsg
		response.Data = nil
	} else {
		response.Error = nil
		response.Data = out
	}

	jsonResponse, jsonErr := json.Marshal(response)
	if jsonErr != nil {
		// Fallback to simple error format if JSON marshaling fails
		return jsonErr.Error()
	}

	return string(jsonResponse)
}

func (a *App) Recompose(shards []string) string {
	out, err := shamir.Recompose(shards)

	response := Response{}
	if err != nil {
		errorMsg := err.Error()
		response.Error = &errorMsg
		response.Data = nil
	} else {
		response.Error = nil
		response.Data = string(out)
	}

	jsonResponse, jsonErr := json.Marshal(response)
	if jsonErr != nil {
		// Fallback to simple error format if JSON marshaling fails
		return jsonErr.Error()
	}

	return string(jsonResponse)
}
