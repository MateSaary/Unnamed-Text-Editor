package handling

import "gioui.org/x/markdown"

// Renderer wraps markdown renderer logic.
type Renderer struct {
	*markdown.Renderer
}

// NewRenderer initializes a new markdown renderer.
func NewRenderer() *Renderer {
	r := markdown.NewRenderer()
	r.Config.MonospaceFont.Typeface = "Go Mono"
	return &Renderer{Renderer: r}
}
