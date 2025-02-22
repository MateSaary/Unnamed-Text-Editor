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
	// Theme allows to customize theme, such as font size.
	Theme *Theme
	// CharacterLabel & LineLabel creates labels for the respective counters.
	CharacterLabel *widget.Label
	LineLabel      *widget.Label

	// Search/Replace Sidebar
	// SearchAreaContainer Holds Search UI
	SearchAreaContainer *fyne.Container
	// SearchTermEntry where you can type text to find.
	SearchTermEntry *widget.Entry
	// ReplaceTermEntry where you can type text to replace matched occurrences.
	ReplaceTermEntry *widget.Entry
	// SearchResults displays number of matches.
	SearchResults *widget.Label
	// Matches hold indices of all occurrences.
	Matches []int
	// CurrentMatchIdx keeps track of current match.
	CurrentMatchIdx int
	// SidebarVisible indicates whether sidebar is currently visible.
	SidebarVisible bool
	// OriginalText stores original text before search markers are added.
	OriginalText string

	// Markdown visibility toggle
	ShowMarkdown bool

	ZoomLabel *widget.Label
}

// NewUI initializes the UI.
func NewUI(app fyne.App, win fyne.Window) *UI {
	theme := NewTheme(app)

	ui := &UI{
		App:                 app,
		Window:              win,
		Editor:              widget.NewMultiLineEntry(),
		Markdown:            widget.NewRichTextFromMarkdown(""),
		Theme:               theme,
		CharacterLabel:      widget.NewLabelWithStyle("Characters: 0", fyne.TextAlignLeading, fyne.TextStyle{Bold: false}),
		LineLabel:           widget.NewLabelWithStyle("Lines: 0", fyne.TextAlignLeading, fyne.TextStyle{Bold: false}),
		SearchAreaContainer: container.NewVBox(),
		SearchTermEntry:     widget.NewEntry(),
		ReplaceTermEntry:    widget.NewEntry(),
		SearchResults:       widget.NewLabel("Results: 0"),
		SidebarVisible:      false,
		Matches:             []int{},
		CurrentMatchIdx:     -1,
		OriginalText:        "",
		ShowMarkdown:        true,
		ZoomLabel:           widget.NewLabelWithStyle("ZoomL 100%", fyne.TextAlignCenter, fyne.TextStyle{Bold: false}),
	}

	ui.MenuBar = ui.CreateMenuBar()
	ui.Theme.ApplyTheme()
	ApplyUserTheme(ui)
	ui.Window.Content().Refresh()

	// Update Markdown Preview whenever text changes.
	ui.Editor.OnChanged = func(content string) {
		ui.RenderMarkdown(content)
		ui.UpdateCounts(content)
	}

	return ui
}

// Updates Markdown Preview.
func (ui *UI) RenderMarkdown(input string) {
	ui.Markdown.ParseMarkdown(input)
	ui.Markdown.Refresh()
}

// Zoom In/Out.
func (ui *UI) ZoomIn() {
	ui.Theme.ZoomIn(ui)
	ui.UpdateZoomLabel()
}

func (ui *UI) ZoomOut() {
	ui.Theme.ZoomOut(ui)
	ui.UpdateZoomLabel()
}

// Reset Zoom.
func (ui *UI) ResetZoom() {
	ui.Theme.ZoomPercent = 100
	ui.Theme.ApplyTheme()
	ui.UpdateZoomLabel()
	ui.Window.Content().Refresh()
}

// Toggle visibility of Markdown preview.
func (ui *UI) toggleMarkdownPreview() {
	ui.ShowMarkdown = !ui.ShowMarkdown
	ui.UpdateLayout()
}

// Update character & line counts.
func (ui *UI) UpdateCounts(content string) {
	charCount := len(content)
	lineCount := len(widget.NewTextGridFromString(content).Rows)

	// Update the labels.
	ui.CharacterLabel.SetText(fmt.Sprintf("Characters: %d", charCount))
	ui.LineLabel.SetText(fmt.Sprintf("Lines: %d", lineCount))
}

func (ui *UI) UpdateZoomLabel() {
	ui.ZoomLabel.SetText(fmt.Sprintf("Zoom: %d%%", ui.Theme.ZoomPercent))
	ui.Window.Content().Refresh()
}
