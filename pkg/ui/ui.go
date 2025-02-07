package ui

import (
	"image"
	"image/color"
	"log"

	"gioui.org/app"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"gioui.org/x/component"
	"gioui.org/x/richtext"
	"github.com/inkeliz/giohyperlink"

	handling "github.com/Leda-Editor/Leda-Text-Editor/pkg/handling"
)

type (
	C = layout.Context
	D = layout.Dimensions
)

// UI specifies the user interface.
type UI struct {
	// External systems.
	// Window provides access to the OS window.
	Window    *app.Window
	// Theme contains semantic style data. Extends `material.Theme`.
	Theme     *Theme
	// Renderer tranforms raw text containing markdown into richtext.
	Renderer  *handling.Renderer
	// Core state.
	// Editor retains raw text in an edit buffer.
	Editor    widget.Editor
	// TextState retains rich text interactions: clicks, hovers and longpresses.
	TextState richtext.InteractiveText
	// Resize state retains the split between the editor and the rendered text.
	Resize    component.Resize
}

// NewUI initializes the UI.
func NewUI(th *Theme) *UI {
	return &UI{
		Window:   new(app.Window),
		Renderer: handling.NewRenderer(),
		Theme:    th,
		Resize:   component.Resize{Ratio: 0.5},
	}
}

// Loop drives the UI until the window is destroyed.
func (ui *UI) Loop() error {
	var ops op.Ops
	for {
		e := ui.Window.Event()
		giohyperlink.ListenEvents(e)
		switch e := e.(type) {
		case app.DestroyEvent:
			return e.Err
		case app.FrameEvent:
			gtx := app.NewContext(&ops, e)
			ui.Layout(gtx)
			e.Frame(gtx.Ops)
		}
	}
}

// Update processes events from the previous frame, updating state accordingly.
func (ui *UI) Update(gtx C) {
	for {
		o, event, ok := ui.TextState.Update(gtx)
		if !ok {
			break
		}
		switch event.Type {
		case richtext.Click:
			if url, ok := o.Get("url").(string); ok && url != "" {
				if err := giohyperlink.Open(url); err != nil {
					log.Printf("error: opening hyperlink: %v", err)
				}
			}
		case richtext.LongPress:
			log.Println("longpress")
			if url, ok := o.Get("url").(string); ok && url != "" {
				ui.Window.Option(app.Title(url))
			}
		}
	}

	// Handle text changes
	for {
		event, ok := ui.Editor.Update(gtx)
		if !ok {
			break
		}
		if _, ok := event.(widget.ChangeEvent); ok {
			var err error
			ui.Theme.Cache, err = ui.Renderer.Render([]byte(ui.Editor.Text()))
			if err != nil {
				log.Printf("error: rendering markdown: %v", err)
			}
		}
	}
}

// Layout renders the current frame.
func (ui *UI) Layout(gtx C) D {
	ui.Update(gtx)
	return ui.Resize.Layout(gtx,
		func(gtx C) D {
			return layout.UniformInset(unit.Dp(4)).Layout(gtx, func(gtx C) D {
				return material.Editor(ui.Theme.Base, &ui.Editor, "markdown").Layout(gtx)
			})
		},
		func(gtx C) D {
			return layout.UniformInset(unit.Dp(4)).Layout(gtx, func(gtx C) D {
				return richtext.Text(&ui.TextState, ui.Theme.Base.Shaper, ui.Theme.Cache...).Layout(gtx)
			})
		},
		func(gtx C) D {
			rect := image.Rectangle{
				Max: image.Point{
					X: (gtx.Dp(unit.Dp(4))),
					Y: (gtx.Constraints.Max.Y),
				},
			}
			paint.FillShape(gtx.Ops, color.NRGBA{A: 200}, clip.Rect(rect).Op())
			return D{Size: rect.Max}
		},
	)
}
