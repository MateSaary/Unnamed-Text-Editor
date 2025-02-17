package ui

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

// Theme struct extends Fyne's theme with custom font sizes for zoom in/out.
type Theme struct {
	Base     fyne.Theme
	FontSize float32
}

// NewTheme instantiates a theme.
func NewTheme() *Theme {
	return &Theme{
		Base:     theme.DefaultTheme(),
		FontSize: theme.TextSize(),
	}
}

// ZoomIn increases font size.
func (th *Theme) ZoomIn() {
	if th.FontSize < 36 { // Max limit to prevent excessive zooming
		th.FontSize += 2
	}
}

// ZoomOut decreases font size.
func (th *Theme) ZoomOut() {
	if th.FontSize > 8 { // Min limit to keep text readable
		th.FontSize -= 2
	}
}

// TextSize returns custom font size.
func (th *Theme) TextSize() float32 {
	return th.FontSize
}

// ApplyTheme applies custom theme globally.
func (th *Theme) ApplyTheme(app fyne.App) {
	app.Settings().SetTheme(th)
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
