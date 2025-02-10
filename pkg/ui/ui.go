package ui

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// UI specifies the user interface.
type UI struct {
	// External systems.
	// Reference to the Fyne Application.
	App fyne.App
	// Window provides access to the OS window.
	Window fyne.Window
	// Core state.
	// Editor retains raw text in an edit buffer.
	Editor *widget.Entry
	// Markdown retains rich text interactions: clicks, hovers and longpresses.
	Markdown *widget.RichText
	// MenuBar adds a menu to the window.
	MenuBar *fyne.Container
	// CHaracterLabel & LineLabel creates labels for the respective counters.
	CharacterLabel *widget.Label
	LineLabel      *widget.Label
}

// NewUI initializes the UI.
func NewUI(app fyne.App, win fyne.Window) *UI {
	ui := &UI{
		App:            app,
		Window:         win,
		Editor:         widget.NewMultiLineEntry(),
		Markdown:       widget.NewRichTextFromMarkdown(""),
		CharacterLabel: widget.NewLabelWithStyle("Characters: 0", fyne.TextAlignLeading, fyne.TextStyle{Bold: false}),
		LineLabel:      widget.NewLabelWithStyle("Lines: 0", fyne.TextAlignLeading, fyne.TextStyle{Bold: false}),
	}

	ui.MenuBar = ui.CreateMenuBar()

	// Update Markdown Preview whenever text changes.
	ui.Editor.OnChanged = func(content string) {
		ui.RenderMarkdown(content)
		ui.UpdateCounts(content)
	}

	return ui
}

// Updates Markdown Preview with proper markdown formatting.
func (ui *UI) RenderMarkdown(input string) {
	ui.Markdown.ParseMarkdown(input)
	ui.Markdown.Refresh()
}

// Creates a split view with the editor and markdown preview.
func (ui *UI) Layout() fyne.CanvasObject {

	statusBar := container.NewHBox(
		ui.CharacterLabel,
		widget.NewLabel(" | "),
		ui.LineLabel,
	)

	statusBar.Resize(fyne.NewSize(statusBar.MinSize().Width, 20))

	split := container.NewHSplit(
		container.NewScroll(ui.Editor),
		container.NewScroll(ui.Markdown),
	)

	split.SetOffset(0.5)

	return container.NewBorder(nil, statusBar, nil, nil, split)
}

// Creates a functional menu bar.
func (ui *UI) CreateMenuBar() *fyne.Container {
	fileMenu := fyne.NewMenu("File",
		fyne.NewMenuItem("Open", func() {}),
		fyne.NewMenuItem("Save", func() {}),
		fyne.NewMenuItem("Exit", func() { ui.Window.Close() }),
	)

	viewMenu := fyne.NewMenu("View",
		fyne.NewMenuItem("Zoom Out", func() {}),
		fyne.NewMenuItem("Zoom In", func() {}),
		fyne.NewMenuItem("Toggle Dark Mode", func() {}),
	)

	helpMenu := fyne.NewMenu("Help",
		fyne.NewMenuItem("About", func() {}),
	)

	mainMenu := fyne.NewMainMenu(fileMenu, viewMenu, helpMenu)
	ui.Window.SetMainMenu(mainMenu)

	return container.NewVBox()
}

func (ui *UI) UpdateCounts(content string) {
	charCount := len(content)                                    // Count characters
	lineCount := len(widget.NewTextGridFromString(content).Rows) // Count lines

	// Update the labels
	ui.CharacterLabel.SetText(fmt.Sprintf("Characters: %d", charCount))
	ui.LineLabel.SetText(fmt.Sprintf("Lines: %d", lineCount))
}
