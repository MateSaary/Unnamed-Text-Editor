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
)

// OpenThemePickerModal opens a modal window with color pickers and RGB sliders
func OpenThemePickerModal(app fyne.App, win fyne.Window, ui *UI) {
	// Get the correct background color from the current theme
	variant := app.Settings().ThemeVariant()
	bgColor := app.Settings().Theme().Color(theme.ColorNameBackground, variant)

	// Convert theme color to RGBA
	bgRGBA := colorToRGBA(bgColor)

	// Create modal background
	modalBackground := canvas.NewRectangle(bgRGBA)

	// Get current window size
	windowSize := win.Canvas().Size()

	// Get the current font size to adjust modal size dynamically
	currentFontSize := ui.Theme.FontSize
	// Set modal size
	modalWidth := windowSize.Width * 0.5
	modalHeight := windowSize.Height * 0.8

	// Adjust modal size based on zoom level
	minHeight := 200 + int(currentFontSize*3) // Increase height slightly based on zoom
	maxHeight := int(windowSize.Height * 0.8) // Max modal height (60% of window height)

	// Ensure modal stays within the min and max height limits
	if modalHeight < float32(minHeight) {
		modalHeight = float32(minHeight)
	}
	if modalHeight > float32(maxHeight) {
		modalHeight = float32(maxHeight)
	}

	// Ensure background fully covers modal
	modalBackground.FillColor = bgRGBA
	modalBackground.Refresh()

	// Retrieve saved theme colors (or fallback to current theme)
	bg := loadColor(app, custom_bg, colorToRGBA(theme.Color(theme.ColorNameBackground)))
	fg := loadColor(app, custom_fg, colorToRGBA(theme.Color(theme.ColorNameForeground)))
	primary := loadColor(app, custom_primary, colorToRGBA(theme.Color(theme.ColorNamePrimary)))
	editorBg := loadColor(app, custom_editor_bg, colorToRGBA(theme.Color(theme.ColorNameInputBackground)))
	menuBg := loadColor(app, custom_menu_bg, colorToRGBA(theme.Color(theme.ColorNameMenuBackground)))

	// Ensure modal background correctly matches selected background
	modalBackground.FillColor = bg
	modalBackground.Refresh()

	// Color preview rectangles
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

	// Declare modal
	var modal dialog.Dialog

	// Close button
	closeBtn := widget.NewButton("Close", func() {
		modal.Hide() // Hide modal when clicked
	})

	// Reset button
	resetThemeBtn := widget.NewButton("Reset", func() {
		ResetCustomTheme(app, ui)
		modal.Hide()
	})

	// Save button
	saveThemeBtn := widget.NewButton("Save", func() {
		SetCustomTheme(app, bg, fg, primary, editorBg, menuBg)
		ui.Window.SetContent(ui.ApplyThemeToLayout())
		ui.Window.Content().Refresh()
		modal.Hide()
	})

	// Space out Save & Close buttons
	buttonsContainer := container.NewHBox(
		layout.NewSpacer(),
		saveThemeBtn,
		layout.NewSpacer(),
		resetThemeBtn,
		layout.NewSpacer(),
		closeBtn,
		layout.NewSpacer(),
	)

	// Create modal layout
	modalContent := container.NewStack(
		modalBackground,
		container.NewVBox(
			createRGBSliders("Background Colour", &bg, bgPreview),
			createRGBSliders("Text Colour", &fg, fgPreview),
			createRGBSliders("Button Colour", &primary, primaryPreview),
			createRGBSliders("Editor Background Colour", &editorBg, editorBgPreview),
			createRGBSliders("Menu Background Colour", &menuBg, menuBgPreview),
			buttonsContainer,
		),
	)

	scrollContainer := container.NewVScroll(modalContent) // Make modal scrollable
	modal = dialog.NewCustomWithoutButtons("Custom Theme Picker", scrollContainer, win)

	// Set the initial modal size
	modal.Resize(fyne.NewSize(modalWidth, modalHeight))

	// Ensure dynamic resizing while keeping the modal centered
	go func() {
		for {
			time.Sleep(time.Millisecond * 200)
			newSize := win.Canvas().Size()

			// Adjust modal size based on zoom level
			modalWidth = newSize.Width * 0.5
			modalHeight = newSize.Height * 0.8

			// Apply constraints to prevent cutting off
			if modalHeight < float32(minHeight) {
				modalHeight = float32(minHeight)
			}
			if modalHeight > float32(maxHeight) {
				modalHeight = float32(maxHeight)
			}

			// Resize modal
			modal.Resize(fyne.NewSize(modalWidth, modalHeight))

		}
	}()

	// Show modal
	modal.Show()
}
