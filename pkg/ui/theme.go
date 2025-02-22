package ui

import (
	"fmt"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
)


// Theme struct extends Fyne's theme with custom font sizes for zoom in/out.
type Theme struct {
	App      fyne.App // Add this field to access app preferences
	Base     fyne.Theme
	FontSize float32
}

const (
	theme_variant    = "theme_variant"
	custom_bg        = "custom_bg"
	custom_fg        = "custom_fg"
	custom_primary   = "custom_primary"
	custom_editor_bg = "custom_editor_bg"
	custom_menu_bg   = "custom_menu_bg"
	custom_button_bg = "custom_button_bg"
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
	currentFontSize := ui.Theme.FontSize

	// If a custom theme is active, reset to default before toggling
	if savedTheme == "custom" {
		resetToDefaultTheme(app) // Fully resets preferences
		savedTheme = "light"     // Default to light mode after reset
	}

	// Toggle between light and dark
	var newTheme fyne.Theme
	if savedTheme == "dark" {
		app.Preferences().SetString(theme_variant, "light")
		newTheme = theme.LightTheme()
	} else {
		app.Preferences().SetString(theme_variant, "dark")
		newTheme = theme.DarkTheme()
	}

	ui.Theme = &Theme{
		App:      app,
		Base:     newTheme,
		FontSize: currentFontSize,
	}

	// Apply the new theme to the UI
	app.Settings().SetTheme(ui.Theme)
	ApplyUserTheme(ui)            // Restores UI with correct theme
	ui.Window.Content().Refresh() // Forces refresh
}

// Resets the theme back to the default Fyne theme and removes saved colors.
func resetToDefaultTheme(app fyne.App) {
	// Clear all custom colors from storage
	app.Preferences().RemoveValue(custom_bg)
	app.Preferences().RemoveValue(custom_fg)
	app.Preferences().RemoveValue(custom_primary)
	app.Preferences().RemoveValue(custom_editor_bg)
	app.Preferences().RemoveValue(custom_menu_bg)
	app.Preferences().RemoveValue(custom_button_bg)

	app.Settings().SetTheme(&Theme{
		App:      app,
		Base:     theme.DefaultTheme(),
		FontSize: theme.TextSize(), // Preserve font size
	})
}

func (ui *UI) ApplyThemeToLayout() fyne.CanvasObject {
	// Create a background rectangle based on the current theme
	bgColor := ui.App.Settings().Theme().Color(theme.ColorNameBackground, ui.App.Settings().ThemeVariant())
	backgroundRect := canvas.NewRectangle(bgColor)

	// Stack the background and actual UI layout
	return container.NewStack(
		backgroundRect, // Background color layer
		ui.Layout(),    // Existing UI content
	)
}

// SetCustomTheme allows switching to a fully custom theme.
func SetCustomTheme(app fyne.App, bg, fg, primary, editorBg, menuBg color.Color) {
	var currentFontSize float32 = theme.TextSize()

	// If the current theme is of type *Theme, keep its FontSize
	if existingTheme, ok := app.Settings().Theme().(*Theme); ok {
		currentFontSize = existingTheme.FontSize
	}

	buttonBg := primary

	// Apply the new theme but **keep zoom level**
	newTheme := &Theme{
		App:      app,
		Base:     NewCustomTheme(bg, fg, primary, editorBg, menuBg, buttonBg),
		FontSize: currentFontSize, // Preserve zoom level
	}

	app.Settings().SetTheme(newTheme)
	app.Preferences().SetString(theme_variant, "custom")

	// Save the colors in Preferences
	saveColor(app, custom_bg, bg)
	saveColor(app, custom_fg, fg)
	saveColor(app, custom_primary, primary)
	saveColor(app, custom_editor_bg, editorBg)
	saveColor(app, custom_menu_bg, menuBg)
	saveColor(app, custom_button_bg, buttonBg)
}

// ResetCustomTheme removes all custom colors and resets to default theme.
func ResetCustomTheme(app fyne.App, ui *UI) {
	currentFontSize := ui.Theme.FontSize
	// Remove all custom theme preferences
	app.Preferences().RemoveValue(custom_bg)
	app.Preferences().RemoveValue(custom_fg)
	app.Preferences().RemoveValue(custom_primary)
	app.Preferences().RemoveValue(custom_editor_bg)
	app.Preferences().RemoveValue(custom_menu_bg)
	app.Preferences().RemoveValue(custom_button_bg)

	// Reset to default Fyne theme
	app.Settings().SetTheme(theme.DefaultTheme())

	// Restore UI with default theme
	ui.Theme = &Theme{
		App:      app,
		Base:     theme.DefaultTheme(),
		FontSize: currentFontSize, 
	}

	app.Settings().SetTheme(ui.Theme)
	ApplyUserTheme(ui)
	ui.Window.Content().Refresh()
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

// NewTheme instantiates a theme.
func NewTheme(app fyne.App) *Theme {
	return &Theme{
		App:      app,
		Base:     theme.DefaultTheme(),
		FontSize: theme.TextSize(),
	}
}

// ZoomIn increases font size.
func (th *Theme) ZoomIn() {
	if th.FontSize < 36 { // Max limit to prevent excessive zooming
		th.FontSize += 2
	}
	th.ApplyTheme()
}

// ZoomOut decreases font size.
func (th *Theme) ZoomOut() {
	if th.FontSize > 8 { // Min limit to keep text readable
		th.FontSize -= 2
	}
	th.ApplyTheme()
}

// TextSize returns custom font size.
func (th *Theme) TextSize() float32 {
	return th.FontSize
}

// ApplyTheme applies the theme while preserving custom colors.
func (th *Theme) ApplyTheme() {
	app := fyne.CurrentApp() // Get the current Fyne app instance

	// Load custom colors
	bg := loadColor(app, custom_bg, colorToRGBA(theme.Color(theme.ColorNameBackground)))
	fg := loadColor(app, custom_fg, colorToRGBA(theme.Color(theme.ColorNameForeground)))
	primary := loadColor(app, custom_primary, colorToRGBA(theme.Color(theme.ColorNamePrimary)))
	editorBg := loadColor(app, custom_editor_bg, colorToRGBA(theme.Color(theme.ColorNameInputBackground)))
	menuBg := loadColor(app, custom_menu_bg, colorToRGBA(theme.Color(theme.ColorNameMenuBackground)))
	buttonBg := loadColor(app, custom_button_bg, primary)

	// Apply the new theme while maintaining colors and font size
	newTheme := &Theme{
		Base:     NewCustomTheme(bg, fg, primary, editorBg, menuBg, buttonBg),
		FontSize: th.FontSize,
	}

	app.Settings().SetTheme(newTheme)
}

// Size returns the custom (or base) size.
func (th *Theme) Size(name fyne.ThemeSizeName) float32 {
	if name == theme.SizeNameText {
		return th.FontSize
	}
	return th.Base.Size(name)
}

// The following is needed as fyne's theme requires an implementation for them. They use base/default implementations.
// Color returns the default color.
func (th *Theme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	return th.Base.Color(name, variant)
}

// Font returns the default font resource.
func (th *Theme) Font(style fyne.TextStyle) fyne.Resource {
	return th.Base.Font(style)
}

// Icon returns the default icon resource.
func (th *Theme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return th.Base.Icon(name)
}
