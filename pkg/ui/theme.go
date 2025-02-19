package ui

import (
	"fmt"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
)

const (
	theme_variant    = "theme_variant"
	custom_bg        = "custom_bg"
	custom_fg        = "custom_fg"
	custom_primary   = "custom_primary"
	custom_editor_bg = "custom_editor_bg"
	custom_menu_bg   = "custom_menu_bg"
)

// ApplyUserTheme applies the last saved theme (light, dark, or custom).
func ApplyUserTheme(ui *UI) {
	// Get the correct background color from the active theme
	variant := ui.App.Settings().ThemeVariant() // Detect current theme mode
	bgColor := ui.App.Settings().Theme().Color(theme.ColorNameBackground, variant)

	// Create a rectangle with the correct background color
	backgroundRect := canvas.NewRectangle(bgColor)

	// Apply the background color and reapply UI layout
	ui.Window.SetContent(container.NewStack(
		backgroundRect, // Background color layer
		ui.Layout(),    // Existing UI content
	))

	// Ensure the window updates
	ui.Window.Content().Refresh()
}

// ToggleDarkMode switches between light and dark themes dynamically.
func ToggleDarkMode(app fyne.App, ui *UI) {
	// Retrieve saved theme preference
	savedTheme := app.Preferences().StringWithFallback(theme_variant, "default")

	// If a custom theme is active, reset to default before toggling
	if savedTheme == "custom" {
		resetToDefaultTheme(app) // Fully resets preferences
		savedTheme = "light"     // Default to light mode after reset
	}

	// Toggle between light and dark
	if savedTheme == "dark" {
		app.Preferences().SetString(theme_variant, "light")
		app.Settings().SetTheme(theme.LightTheme()) // Switch to light mode
	} else {
		app.Preferences().SetString(theme_variant, "dark")
		app.Settings().SetTheme(theme.DarkTheme()) // Switch to dark mode
	}

	// **Apply the new theme to the UI**
	ui.Window.SetContent(ui.ApplyThemeToLayout()) // Restores UI with correct theme
	ui.Window.Content().Refresh()                 // Forces refresh
}

// Resets the theme back to the default Fyne theme and removes saved colors.
func resetToDefaultTheme(app fyne.App) {
	// Clear all custom colors from storage
	app.Preferences().RemoveValue(custom_bg)
	app.Preferences().RemoveValue(custom_fg)
	app.Preferences().RemoveValue(custom_primary)
	app.Preferences().RemoveValue(custom_editor_bg)
	app.Preferences().RemoveValue(custom_menu_bg)

	// Reset to default Fyne theme
	app.Settings().SetTheme(theme.DefaultTheme())
}

func (ui *UI) ApplyThemeToLayout() fyne.CanvasObject {
	// Create a background rectangle based on the current theme
	bgColor := theme.Color(theme.ColorNameBackground)
	backgroundRect := canvas.NewRectangle(bgColor)

	// Stack the background and actual UI layout
	return container.NewStack(
		backgroundRect, // Background color layer
		ui.Layout(),    // Existing UI content
	)
}

// SetCustomTheme allows switching to a fully custom theme.
func SetCustomTheme(app fyne.App, bg, fg, primary, editorBg, menuBg color.Color) {
	app.Settings().SetTheme(NewCustomTheme(bg, fg, primary, editorBg, menuBg)) // Apply custom theme
	app.Preferences().SetString(theme_variant, "custom")

	// Save the colors in Preferences
	saveColor(app, custom_bg, bg)
	saveColor(app, custom_fg, fg)
	saveColor(app, custom_primary, primary)
	saveColor(app, custom_editor_bg, editorBg)
	saveColor(app, custom_menu_bg, menuBg)
}

// Helper functions to save and load colors
func saveColor(app fyne.App, key string, c color.Color) {
	r, g, b, a := c.RGBA()
	app.Preferences().SetString(key, formatColor(r, g, b, a))
}

func loadColor(app fyne.App, key string, fallback color.RGBA) color.RGBA {
	colorStr := app.Preferences().StringWithFallback(key, formatColorRGBA(fallback))
	return parseColor(colorStr, fallback)
}

// Helper functions to format and parse colors
func formatColor(r, g, b, a uint32) string {
	return fmt.Sprintf("%d,%d,%d,%d", r>>8, g>>8, b>>8, a>>8)
}

func formatColorRGBA(c color.Color) string {
	r, g, b, a := c.RGBA()
	return formatColor(r, g, b, a)
}

func parseColor(s string, fallback color.RGBA) color.RGBA {
	var r, g, b, a int
	_, err := fmt.Sscanf(s, "%d,%d,%d,%d", &r, &g, &b, &a)
	if err != nil {
		fyne.LogError("Failed to parse color", err)
		return fallback // Return fallback if parsing fails
	}
	return color.RGBA{uint8(r), uint8(g), uint8(b), uint8(a)}
}

func colorToRGBA(c color.Color) color.RGBA {
	r, g, b, a := c.RGBA()
	return color.RGBA{uint8(r >> 8), uint8(g >> 8), uint8(b >> 8), uint8(a >> 8)}
}
