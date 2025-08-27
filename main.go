package main

import (
	"embed"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/mac"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	// Create an instance of the app structure
	app := NewApp()

	// Create application with options
	err := wails.Run(&options.App{
		Title:     "orcrux",
		Width:     1024,
		Height:    768,
		MinWidth:  1024,
		MinHeight: 768,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		// todo: set bg to dark green
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 27, A: 1},
		OnStartup:        app.startup,
		Bind: []interface{}{
			app,
		},
		Mac: &mac.Options{
			TitleBar: &mac.TitleBar{
				HideTitle:                  false,
				TitlebarAppearsTransparent: false,
				UseToolbar:                 false,
			},
			About: &mac.AboutInfo{
				Title:   "orcrux",
				Message: "© 2025 Luc T",
			},
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
