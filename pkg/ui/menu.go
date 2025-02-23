package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	handling "github.com/Leda-Editor/Leda-Text-Editor/pkg/handling"
)

// Creates a functional menu bar.
func (ui *UI) CreateMenuBar() *fyne.Container {
	fileMenu := fyne.NewMenu("File",
		fyne.NewMenuItem("Open", func() { handling.OpenFile(ui.Window, ui.Editor) }),
		fyne.NewMenuItem("Save", func() { handling.SaveFile(ui.Window, ui.Editor) }),
		fyne.NewMenuItem("Exit", func() { handling.ClearEditor(ui.Editor) }),
	)

	viewMenu := fyne.NewMenu("View",
		fyne.NewMenuItem("Zoom Out", func() { ui.ZoomOut() }),
		fyne.NewMenuItem("Zoom In", func() { ui.ZoomIn() }),
		fyne.NewMenuItem("Reset Zoom", func() { ui.ResetZoom() }),
		fyne.NewMenuItemSeparator(),
		fyne.NewMenuItem("Show/Hide Markdown Preview", func() { ui.toggleMarkdownPreview() }),
		fyne.NewMenuItem("Dark Mode On/Off", func() { ToggleDarkMode(ui.App, ui) }),
		fyne.NewMenuItem("Set Custom Theme", func() {
			OpenThemePickerModal(ui.App, ui.Window, ui)
		}),
	)

	editMenu := fyne.NewMenu("Edit",
		fyne.NewMenuItem("Search", func() { ShowSearchUI(false, ui) }),
		fyne.NewMenuItem("Search & Replace", func() { ShowSearchUI(true, ui) }),
	)

	helpMenu := fyne.NewMenu("Help",
		fyne.NewMenuItem("About", func() {
			dialog.ShowInformation("About", "\n\nLeda is a text editor built with Go and Fyne.", ui.Window)
		}),
	)

	mainMenu := fyne.NewMainMenu(fileMenu, viewMenu, editMenu, helpMenu)
	ui.Window.SetMainMenu(mainMenu)

	return container.NewVBox()
}
