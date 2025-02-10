package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	ui "github.com/Leda-Editor/Leda-Text-Editor/pkg/ui"
)

func main() {
	// Initialize Fyne Application.
	a := app.New()

	// Create a new window for the application
	w := a.NewWindow("Leda Text Editor")

	// Initialize UI.
	ledaUI := ui.NewUI(a, w)

	// Set up window layout.
	w.SetContent(ledaUI.Layout())

	// Set window size.
	w.Resize(fyne.NewSize(800, 600))
	// Display the window and start the event loop.
	w.ShowAndRun()
}
