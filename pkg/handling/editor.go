package handling

import "gioui.org/widget"

// Editor wraps the text editor widget.
type Editor struct {
	widget.Editor
}

// NewEditor creates a new markdown editor.
func NewEditor() *Editor {
	return &Editor{}
}
