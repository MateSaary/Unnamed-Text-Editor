package main

import (
	"log"
	"os"

	"gioui.org/app"
	"gioui.org/font/gofont"

	ui "github.com/Leda-Editor/Leda-Text-Editor/pkg/ui"
)

func main() {
	th := ui.NewTheme(gofont.Collection())
	ui := ui.NewUI(th)

	go func() {
		if err := ui.Loop(); err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()

	app.Main()
}
