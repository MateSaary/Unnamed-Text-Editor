package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	ui "github.com/Leda-Editor/Leda-Text-Editor/pkg/ui"
)

func main() {
	// Initialize Fyne Application.
	app := app.NewWithID("leda-text-editor")

	// Create a new window for the application
	window := app.NewWindow("Leda Text Editor")

	// Initialize UI.
	ledaUI := ui.NewUI(app, window)

	// Set up window layout.
	window.SetContent(ledaUI.Layout())

	// Set window size.
	window.Resize(fyne.NewSize(900, 700))
	// Display the window and start the event loop.
	window.ShowAndRun()
}
