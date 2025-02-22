package ui

import (
	"fmt"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
)

// Perform search and highlight results.
func (ui *UI) performSearch() {
	term := ui.SearchTermEntry.Text
	if term == "" {
		ui.SearchResults.SetText("Results: 0")
		ui.Matches = []int{}
		ui.RenderMarkdown(ui.Editor.Text)
		return
	}

	if ui.OriginalText == "" {
		ui.OriginalText = ui.Editor.Text
	}

	ui.Editor.SetText(ui.OriginalText)

	ui.Matches = []int{}
	text := ui.Editor.Text

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

// Scroll to specific match.
func (ui *UI) scrollToMatch(idx int) {
	if len(ui.Matches) == 0 || idx < 0 || idx >= len(ui.Matches) {
		return
	}

	ui.Editor.SetText(ui.OriginalText)

	// Get the current match index.
	matchIdx := ui.Matches[idx]
	ui.CurrentMatchIdx = idx
	// Insert arrows around the match.
	searchTerm := ui.SearchTermEntry.Text
	matchLen := len(searchTerm)

	if matchIdx < 0 || matchIdx+matchLen > len(ui.OriginalText) {
		return
	}

	highlighted := ui.OriginalText[:matchIdx] + "⬅️" + searchTerm + "➡️" + ui.OriginalText[matchIdx+matchLen:]
	ui.Editor.SetText(highlighted)

	ui.Editor.CursorColumn = matchIdx
	ui.Editor.CursorRow = strings.Count(ui.Editor.Text[:matchIdx], "\n")

	ui.ensureVisible(matchIdx)

	ui.Editor.Refresh()
}

// Ensure the current match is fully visible.
func (ui *UI) ensureVisible(position int) {
	if position < 0 || position >= len(ui.Editor.Text) {
		return
	}

	// Calculate line and column.
	line := strings.Count(ui.Editor.Text[:position], "\n")
	col := position - strings.LastIndex(ui.Editor.Text[:position], "\n") - 1

	ui.Editor.CursorRow = line
	ui.Editor.CursorColumn = col

	visibleWidth := 40
	if col > visibleWidth/2 {
		ui.Editor.CursorColumn = col - (visibleWidth / 2)
	} else {
		ui.Editor.CursorColumn = 0
	}

	if col+len(ui.SearchTermEntry.Text) < len(ui.Editor.Text)-10 {
		ui.Editor.CursorColumn = col + 10
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

	ui.Editor.SetText(ui.OriginalText)
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

	ui.Editor.SetText(ui.OriginalText)
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

	if ui.OriginalText == "" {
		ui.OriginalText = ui.Editor.Text
	}

	currentIdx := ui.Matches[ui.CurrentMatchIdx]

	if currentIdx >= 0 && currentIdx+len(term) <= len(ui.OriginalText) {
		newText := ui.OriginalText[:currentIdx] + replace + ui.OriginalText[currentIdx+len(term):]
		ui.Editor.SetText(newText)
		ui.OriginalText = newText
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

	if ui.OriginalText == "" {
		ui.OriginalText = ui.Editor.Text
	}

	newText := strings.ReplaceAll(ui.OriginalText, term, replace)
	ui.Editor.SetText(newText)
	ui.OriginalText = newText

	ui.performSearch()
}

// Toggle sidebar visibility.
func (ui *UI) toggleSidebar() {
	if ui.SidebarVisible {
		if ui.OriginalText != "" {
			ui.Editor.SetText(ui.OriginalText)
		}
		ui.OriginalText = ""
	}

	ui.SidebarVisible = !ui.SidebarVisible
	ui.UpdateLayout()
}

func ShowSearchContainer(ui *UI, content fyne.CanvasObject) {
	if ui.SearchAreaContainer == nil {
		ui.SearchAreaContainer = container.NewVBox()
	}

	ui.SearchAreaContainer.Objects = []fyne.CanvasObject{content}
	ui.SearchAreaContainer.Refresh()
}

func ShowSearchUI(enableReplace bool, ui *UI) {
	if ui.SearchAreaContainer == nil {
		ui.SearchAreaContainer = container.NewVBox()
	}

	ui.ReplaceTermEntry.Hide()
	if enableReplace {
		ui.ReplaceTermEntry.Show()
	}

	ui.SidebarVisible = true
	ui.UpdateLayout()
}
