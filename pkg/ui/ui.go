package ui

import (
	"fmt"
	"image/color"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
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
	// CharacterLabel & LineLabel creates labels for the respective counters.
	CharacterLabel *widget.Label
	LineLabel      *widget.Label
}

// NewUI initializes the UI.
func NewUI(app fyne.App, win fyne.Window) *UI {
	win.SetPadded(false)

	ui := &UI{
		App:            app,
		Window:         win,
		Editor:         widget.NewMultiLineEntry(),
		Markdown:       widget.NewRichTextFromMarkdown(""),
		CharacterLabel: widget.NewLabelWithStyle("Characters: 0", fyne.TextAlignLeading, fyne.TextStyle{Bold: false}),
		LineLabel:      widget.NewLabelWithStyle("Lines: 0", fyne.TextAlignLeading, fyne.TextStyle{Bold: false}),
	}

	ui.MenuBar = ui.CreateMenuBar()

	ApplyUserTheme(ui)

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

// OpenThemePickerModal opens a modal window with color pickers and RGB sliders
func OpenThemePickerModal(app fyne.App, win fyne.Window, ui *UI) {
	// Get the correct background color from the current theme
	variant := app.Settings().ThemeVariant()
	bgColor := app.Settings().Theme().Color(theme.ColorNameBackground, variant)

	// Convert theme color to RGBA
	bgRGBA := colorToRGBA(bgColor)

	// **Create modal background**
	modalBackground := canvas.NewRectangle(bgRGBA)

	// Get initial window size
	windowSize := win.Canvas().Size()

	// **Set modal size (50% width, 40% height)**
	modalWidth := windowSize.Width * 0.5
	modalHeight := windowSize.Height * 0.4

	// **Ensure modal never becomes too small**
	if modalWidth < 350 {
		modalWidth = 350
	}
	if modalHeight < 200 {
		modalHeight = 200
	}

	// **Ensure background fully covers modal**
	modalBackground.FillColor = bgRGBA
	modalBackground.Refresh()

	// **Retrieve saved theme colors (or fallback to current theme)**
	bg := loadColor(app, "custom_bg", colorToRGBA(theme.Color(theme.ColorNameBackground)))
	fg := loadColor(app, "custom_fg", colorToRGBA(theme.Color(theme.ColorNameForeground)))
	primary := loadColor(app, "custom_primary", colorToRGBA(theme.Color(theme.ColorNamePrimary)))
	editorBg := loadColor(app, "custom_editor_bg", colorToRGBA(theme.Color(theme.ColorNameInputBackground)))
	menuBg := loadColor(app, "custom_menu_bg", colorToRGBA(theme.Color(theme.ColorNameMenuBackground)))

	// **Ensure modal background correctly matches selected background**
	modalBackground.FillColor = bg
	modalBackground.Refresh()

	// **Color preview rectangles**
	bgPreview := canvas.NewRectangle(bg)
	fgPreview := canvas.NewRectangle(fg)
	primaryPreview := canvas.NewRectangle(primary)
	editorBgPreview := canvas.NewRectangle(editorBg)
	menuBgPreview := canvas.NewRectangle(menuBg)

	// Helper function to create RGB sliders
	createRGBSliders := func(label string, c *color.RGBA, preview *canvas.Rectangle) fyne.CanvasObject {
		r := widget.NewSlider(0, 255)
		g := widget.NewSlider(0, 255)
		b := widget.NewSlider(0, 255)

		// Set slider values to current color values
		r.Value, g.Value, b.Value = float64(c.R), float64(c.G), float64(c.B)

		// RGB Values Label
		rgbValueLabel := widget.NewLabel(fmt.Sprintf("R:%d G:%d B:%d", c.R, c.G, c.B))

		// Function to update color preview and label
		updateRGBLabel := func() {
			rgbValueLabel.SetText(fmt.Sprintf("R:%d G:%d B:%d", c.R, c.G, c.B))
			preview.FillColor = *c
			preview.Refresh()
		}

		r.OnChanged = func(v float64) { c.R = uint8(v); updateRGBLabel() }
		g.OnChanged = func(v float64) { c.G = uint8(v); updateRGBLabel() }
		b.OnChanged = func(v float64) { c.B = uint8(v); updateRGBLabel() }

		return container.NewVBox(
			widget.NewLabel(label),
			container.NewGridWithColumns(4, r, g, b, preview),
			rgbValueLabel,
		)
	}

	// **Declare modal**
	var modal dialog.Dialog

	// "Close" button
	closeBtn := widget.NewButton("Close", func() {
		modal.Hide() // Hide modal when clicked
	})

	// "Save Theme" button
	saveThemeBtn := widget.NewButton("Save Theme", func() {
		SetCustomTheme(app, bg, fg, primary, editorBg, menuBg)
		ui.Window.SetContent(ui.ApplyThemeToLayout())
		ui.Window.Content().Refresh()
		modal.Hide()
	})

	// **Space out Save & Close buttons**
	buttonsContainer := container.NewHBox(
		layout.NewSpacer(),
		saveThemeBtn,
		layout.NewSpacer(),
		closeBtn,
		layout.NewSpacer(),
	)

	// **Create modal layout**
	modalContent := container.NewStack(
		modalBackground,
		container.NewVBox(
			createRGBSliders("Background", &bg, bgPreview),
			createRGBSliders("Foreground", &fg, fgPreview),
			createRGBSliders("Primary", &primary, primaryPreview),
			createRGBSliders("Editor Background", &editorBg, editorBgPreview),
			createRGBSliders("Menu Background", &menuBg, menuBgPreview),
			buttonsContainer,
		),
	)

	modal = dialog.NewCustomWithoutButtons("Custom Theme Picker", modalContent, win)
	// **Set the initial modal size**
	modal.Resize(fyne.NewSize(modalWidth, modalHeight))

	// **Ensure dynamic resizing while keeping the modal centered**
	go func() {
		for {
			time.Sleep(time.Millisecond * 200)
			newSize := win.Canvas().Size()

			// Keep modal centered while resizing
			modal.Resize(fyne.NewSize(newSize.Width*0.5, newSize.Height*0.4))

			// Shrink when minimized
			if newSize.Height < 400 {
				modal.Resize(fyne.NewSize(newSize.Width*0.5, newSize.Height*0.3))
			}
		}
	}()

	// **Show modal**
	modal.Show()
}

// Creates a functional menu bar.
func (ui *UI) CreateMenuBar() *fyne.Container {
	fileMenu := fyne.NewMenu("File",
		fyne.NewMenuItem("Open", func() { handling.OpenFile(ui.Window, ui.Editor) }),
		fyne.NewMenuItem("Save", func() { handling.SaveFile(ui.Window, ui.Editor) }),
		fyne.NewMenuItem("Exit", func() { handling.ClearEditor(ui.Editor) }),
	)

	viewMenu := fyne.NewMenu("View",
		fyne.NewMenuItem("Zoom Out", func() {}),
		fyne.NewMenuItem("Zoom In", func() {}),
		fyne.NewMenuItem("Toggle Dark Mode", func() { ToggleDarkMode(ui.App, ui) }),
		fyne.NewMenuItem("Set Custom Theme", func() {
			OpenThemePickerModal(ui.App, ui.Window, ui)
		}),
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
