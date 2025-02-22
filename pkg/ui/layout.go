package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

// Creates a split view with the editor and markdown preview.
func (ui *UI) Layout() fyne.CanvasObject {
	separator := widget.NewSeparator()

	statusBar := container.NewVBox(
		separator,
		container.NewHBox(
			ui.CharacterLabel,
			widget.NewLabel(" | "),
			ui.LineLabel,
			layout.NewSpacer(),
			ui.ZoomLabel,
		),
	)

	var sidebar fyne.CanvasObject = container.NewVBox()

	if ui.SidebarVisible {
		sidebar = container.NewVBox(
			widget.NewLabel("🔍 Search"),
			ui.SearchTermEntry,
			widget.NewButton("Search", func() { ui.performSearch() }),
			ui.SearchResults,
			widget.NewSeparator(),
		)

		if ui.ReplaceTermEntry.Visible() {
			sidebar.(*fyne.Container).Add(widget.NewLabel("🔁 Replace"))
			sidebar.(*fyne.Container).Add(ui.ReplaceTermEntry)
			sidebar.(*fyne.Container).Add(widget.NewButton("Replace Current", func() { ui.performReplaceCurrent() }))
			sidebar.(*fyne.Container).Add(widget.NewButton("Replace All", func() { ui.performReplaceAll() }))
			sidebar.(*fyne.Container).Add(widget.NewSeparator())
		}

		sidebar.(*fyne.Container).Add(widget.NewButton("⬆️ Previous", func() { ui.previousMatch() }))
		sidebar.(*fyne.Container).Add(widget.NewButton("⬇️ Next", func() { ui.nextMatch() }))
		sidebar.(*fyne.Container).Add(widget.NewButton("❌ Close", func() { ui.toggleSidebar() }))
	}

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

func (ui *UI) UpdateLayout() {
	ui.Window.SetContent(ui.Layout())
	ui.Window.Content().Refresh()
}
