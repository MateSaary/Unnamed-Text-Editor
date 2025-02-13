package ui

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

// CustomTheme defines a theme with user-defined colors.
type CustomTheme struct {
	Background color.Color
	Foreground color.Color
	Primary    color.Color
	EditorBg   color.Color
	MenuBg     color.Color
}

// Color overrides default colors with custom values.
func (t *CustomTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	switch name {
	case theme.ColorNameBackground:
		return t.Background
	case theme.ColorNameForeground:
		return t.Foreground
	case theme.ColorNamePrimary:
		return t.Primary
	case theme.ColorNameInputBackground:
		return t.EditorBg
	case theme.ColorNameMenuBackground:
		return t.MenuBg
	}
	return theme.DefaultTheme().Color(name, variant) // Fallback to default
}

// Font overrides the default font.
func (t *CustomTheme) Font(style fyne.TextStyle) fyne.Resource {
	return theme.DefaultTheme().Font(style) // Keep default font
}

// Icon overrides the default icons.
func (t *CustomTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(name) // Keep default icons
}

// Size overrides default sizes.
func (t *CustomTheme) Size(name fyne.ThemeSizeName) float32 {
	return theme.DefaultTheme().Size(name) // Keep default sizes
}

// NewCustomTheme initializes a theme with user-defined colors.
func NewCustomTheme(bg, fg, primary, editorBg, menuBg color.Color) fyne.Theme {
	return &CustomTheme{
		Background: bg,
		Foreground: fg,
		Primary:    primary,
		EditorBg:   editorBg,
		MenuBg:     menuBg,
	}
}
