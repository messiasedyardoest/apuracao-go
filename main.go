package main

import (
	"embed"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

// Saida do `nuxt generate` (SPA, ssr: false).
//
//go:embed all:frontend/.output/public
var assets embed.FS

func main() {
	app := NewApp()

	err := wails.Run(&options.App{
		Title:     "Apuração 2026 · Presidente 2º turno",
		Width:     1366,
		Height:    860,
		MinWidth:  1024,
		MinHeight: 680,
		// Sem moldura do Windows: a barra de titulo e o BarraTitulo.vue.
		Frameless: true,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 241, G: 237, B: 233, A: 1},
		OnStartup:        app.startup,
		OnShutdown:       app.shutdown,
		Bind: []interface{}{
			app,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
