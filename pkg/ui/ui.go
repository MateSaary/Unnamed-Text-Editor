package ui

import (
	"fmt"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	handling "github.com/Leda-Editor/Leda-Text-Editor/pkg/handling"
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
}

// NewUI initializes the UI.
func NewUI(app fyne.App, win fyne.Window) *UI {
	theme := NewTheme()

	ui := &UI{
		App:              app,
		Window:           win,
		Editor:           widget.NewMultiLineEntry(),
		Markdown:         widget.NewRichTextFromMarkdown(""),
		Theme:            theme,
		CharacterLabel:   widget.NewLabelWithStyle("Characters: 0", fyne.TextAlignLeading, fyne.TextStyle{Bold: false}),
		LineLabel:        widget.NewLabelWithStyle("Lines: 0", fyne.TextAlignLeading, fyne.TextStyle{Bold: false}),
		SearchTermEntry:  widget.NewEntry(),
		ReplaceTermEntry: widget.NewEntry(),
		SearchResults:    widget.NewLabel("Results: 0"),
		SidebarVisible:   false,
		Matches:          []int{},
		CurrentMatchIdx:  -1,
		OriginalText:     "",
		ShowMarkdown:     true,
	}

	ui.MenuBar = ui.CreateMenuBar()
	ui.Theme.ApplyTheme(ui.App)

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

// Zoom In/Out.
func (ui *UI) ZoomIn() {
	ui.Theme.ZoomIn()
	ui.Theme.ApplyTheme(ui.App)
}

func (ui *UI) ZoomOut() {
	ui.Theme.ZoomOut()
	ui.Theme.ApplyTheme(ui.App)
}

// Toggle visibility of Markdown preview.
func (ui *UI) toggleMarkdownPreview() {
	ui.ShowMarkdown = !ui.ShowMarkdown
	ui.Window.SetContent(ui.Layout())
}

// Perform search and highlight results.
func (ui *UI) performSearch() {
	term := ui.SearchTermEntry.Text
	if term == "" {
		ui.SearchResults.SetText("Results: 0")
		ui.Matches = []int{}
		ui.RenderMarkdown(ui.Editor.Text)
		return
	}

	ui.Matches = []int{}
	text := ui.Editor.Text
	ui.OriginalText = text

	// Find all occurrences.
	start := 0
	for {
		index := strings.Index(text[start:], term)
		if index == -1 {
			break
		}
		index += start
		ui.Matches = append(ui.Matches, index)
		start = index + len(term)
	}

	count := len(ui.Matches)
	ui.SearchResults.SetText(fmt.Sprintf("Results: %d", count))

	if count == 0 {
		dialog.ShowInformation("Find", fmt.Sprintf("No occurrences of '%s' found.", term), ui.Window)
	} else {
		ui.CurrentMatchIdx = 0
		ui.scrollToMatch(ui.CurrentMatchIdx)
	}
}

// Scroll to specific match .
func (ui *UI) scrollToMatch(idx int) {
	if len(ui.Matches) == 0 {
		return
	}

	cleanedText := strings.ReplaceAll(ui.Editor.Text, "➡️", "")
	cleanedText = strings.ReplaceAll(cleanedText, "⬅️", "")
	ui.Editor.SetText(cleanedText)

	// Get the current match index.
	matchIdx := ui.Matches[idx]
	ui.CurrentMatchIdx = idx

	// Calculate cursor position.
	ui.Editor.CursorColumn = matchIdx
	ui.Editor.CursorRow = strings.Count(ui.Editor.Text[:matchIdx], "\n")

	// Insert arrows around the match.
	searchTerm := ui.SearchTermEntry.Text
	matchLen := len(searchTerm)

	if matchIdx >= 0 && matchIdx+matchLen <= len(ui.Editor.Text) {
		highlighted := ui.Editor.Text[:matchIdx] + "⬅➡️" + searchTerm + "⬅️" + ui.Editor.Text[matchIdx+matchLen:]
		ui.Editor.SetText(highlighted)
	}

	ui.Editor.Refresh()
}

// Navigate to the previous match.
func (ui *UI) previousMatch() {
	if len(ui.Matches) == 0 {
		return
	}
	ui.CurrentMatchIdx--
	if ui.CurrentMatchIdx < 0 {
		ui.CurrentMatchIdx = len(ui.Matches) - 1
	}
	ui.scrollToMatch(ui.CurrentMatchIdx)
}

// Navigate to the next match.
func (ui *UI) nextMatch() {
	if len(ui.Matches) == 0 {
		return
	}
	ui.CurrentMatchIdx++
	if ui.CurrentMatchIdx >= len(ui.Matches) {
		ui.CurrentMatchIdx = 0
	}
	ui.scrollToMatch(ui.CurrentMatchIdx)
}

// Perform replace on current match.
func (ui *UI) performReplaceCurrent() {
	if len(ui.Matches) == 0 || ui.CurrentMatchIdx == -1 {
		dialog.ShowInformation("Replace", "No match selected.", ui.Window)
		return
	}

	term := ui.SearchTermEntry.Text
	replace := ui.ReplaceTermEntry.Text

	// Clean any arrows from the text.
	cleanedText := strings.ReplaceAll(ui.Editor.Text, "➡️", "")
	cleanedText = strings.ReplaceAll(cleanedText, "⬅️", "")

	currentIdx := ui.Matches[ui.CurrentMatchIdx]

	if currentIdx >= 0 && currentIdx+len(term) <= len(cleanedText) {
		newText := cleanedText[:currentIdx] + replace + cleanedText[currentIdx+len(term):]
		ui.Editor.SetText(newText)
	}

	ui.performSearch()
}

// Perform replace-all.
func (ui *UI) performReplaceAll() {
	term := ui.SearchTermEntry.Text
	replace := ui.ReplaceTermEntry.Text

	if term == "" {
		dialog.ShowInformation("Replace", "Enter a search term.", ui.Window)
		return
	}

	cleanedText := strings.ReplaceAll(ui.Editor.Text, "➡️", "")
	cleanedText = strings.ReplaceAll(cleanedText, "⬅️", "")

	if !strings.Contains(cleanedText, term) {
		dialog.ShowInformation("Replace", fmt.Sprintf("'%s' not found.", term), ui.Window)
		return
	}

	ui.Editor.SetText(strings.ReplaceAll(cleanedText, term, replace))
	ui.performSearch()
}

// Toggle sidebar visibility.
func (ui *UI) toggleSidebar() {
	ui.SidebarVisible = !ui.SidebarVisible
	ui.Window.SetContent(ui.Layout())
}

// Creates a split view with the editor and markdown preview.
func (ui *UI) Layout() fyne.CanvasObject {

	statusBar := container.NewHBox(
		ui.CharacterLabel,
		widget.NewLabel(" | "),
		ui.LineLabel,
	)

	sidebar := container.NewVBox(
		widget.NewLabel("🔍 Search"),
		ui.SearchTermEntry,
		widget.NewButton("Search", func() { ui.performSearch() }),
		ui.SearchResults,
		widget.NewSeparator(),
		widget.NewLabel("🔁 Replace"),
		ui.ReplaceTermEntry,
		widget.NewButton("Replace Current", func() { ui.performReplaceCurrent() }),
		widget.NewButton("Replace All", func() { ui.performReplaceAll() }),
		widget.NewSeparator(),
		widget.NewButton("⬆️ Previous", func() { ui.previousMatch() }),
		widget.NewButton("⬇️ Next", func() { ui.nextMatch() }),
		widget.NewButton("❌ Close", func() { ui.toggleSidebar() }),
	)

	var content fyne.CanvasObject
	if ui.SidebarVisible {
		split := container.NewHSplit(
			sidebar,
			container.NewScroll(ui.Editor),
		)
		split.SetOffset(0.2)

		if ui.ShowMarkdown {
			content = container.NewHSplit(
				split,
				container.NewScroll(ui.Markdown),
			)
		} else {
			content = split
		}
	} else {
		if ui.ShowMarkdown {
			content = container.NewHSplit(
				container.NewScroll(ui.Editor),
				container.NewScroll(ui.Markdown),
			)
		} else {
			content = container.NewScroll(ui.Editor)
		}
	}
	return container.NewBorder(nil, statusBar, nil, nil, content)
}

// Creates a functional menu bar.
func (ui *UI) CreateMenuBar() *fyne.Container {
	fileMenu := fyne.NewMenu("File",
		fyne.NewMenuItem("Open", func() { handling.OpenFile(ui.Window, ui.Editor) }),
		fyne.NewMenuItem("Save", func() { handling.SaveFile(ui.Window, ui.Editor) }),
		fyne.NewMenuItem("Exit", func() { handling.ClearEditor(ui.Editor) }),
	)

	viewMenu := fyne.NewMenu("View",
		fyne.NewMenuItem("Zoom In", func() { ui.ZoomIn() }),
		fyne.NewMenuItem("Zoom Out", func() { ui.ZoomOut() }),
		fyne.NewMenuItem("Show/Hide Markdown Preview", func() { ui.toggleMarkdownPreview() }),
		fyne.NewMenuItem("Dark Mode On/Off", func() {}),
	)

	editMenu := fyne.NewMenu("Edit",
		fyne.NewMenuItem("Find/Replace", func() { ui.toggleSidebar() }),
	)

	helpMenu := fyne.NewMenu("Help",
		fyne.NewMenuItem("About", func() {}),
	)

	mainMenu := fyne.NewMainMenu(fileMenu, viewMenu, editMenu, helpMenu)
	ui.Window.SetMainMenu(mainMenu)

	return container.NewVBox()
}

// Update character & line counts.
func (ui *UI) UpdateCounts(content string) {
	charCount := len(content)                                    // Count characters.
	lineCount := len(widget.NewTextGridFromString(content).Rows) // Count lines.

	// Update the labels.
	ui.CharacterLabel.SetText(fmt.Sprintf("Characters: %d", charCount))
	ui.LineLabel.SetText(fmt.Sprintf("Lines: %d", lineCount))
}
