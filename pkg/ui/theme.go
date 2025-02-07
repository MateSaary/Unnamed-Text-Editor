package ui

import (
	"gioui.org/font"
	"gioui.org/text"
	"gioui.org/widget/material"
	"gioui.org/x/richtext"
)

// Theme contains semantic style data.
type Theme struct {
	// Base theme to extend.
	Base  *material.Theme
	// cache of processed markdown.
	Cache []richtext.SpanStyle
}

// NewTheme instantiates a theme, extending material theme.
func NewTheme(fonts []font.FontFace) *Theme {
	th := material.NewTheme()
	th.Shaper = text.NewShaper(text.WithCollection(fonts))
	return &Theme{Base: th}
}
