package ui

import (
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
	// TextState retains rich text interactions: clicks, hovers and longpresses.
	Markdown *widget.RichText
}

// NewUI initializes the UI.
func NewUI(app fyne.App, win fyne.Window) *UI {
	ui := &UI{
		App:    app,
		Window: win,
		// Creates a multi-line text editor.
		Editor: widget.NewMultiLineEntry(),
		// Initializes an empty Markdown preview.
		Markdown: widget.NewRichTextFromMarkdown(""),
	}

	// Update Markdown Preview whenever text changes.
	ui.Editor.OnChanged = func(content string) {
		ui.RenderMarkdown(content)
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
	split := container.NewHSplit(
		container.NewScroll(ui.Editor),
		container.NewScroll(ui.Markdown),
	)

	split.SetOffset(0.5)

	return split
}
