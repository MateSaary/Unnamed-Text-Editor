package main

import (
	"log"
	"os"

	"gioui.org/app"
	"gioui.org/unit"

	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

func main() {

	go func() {
		window := new(app.Window)
		window.Option(app.Title("Leda Text Editor"))

		err := run(window)
		if err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()
	app.Main()

}

func run(window *app.Window) error {
	theme := material.NewTheme()
	var ops op.Ops

	editor := new(widget.Editor)
	editor.SingleLine = false
	for {
		switch e := window.Event().(type) {
		case app.DestroyEvent:
			return e.Err
		case app.FrameEvent:
			gtx := app.NewContext(&ops, e)

			addEditor(&gtx, theme, editor)
			e.Frame(gtx.Ops)
		}
	}
}

func addEditor(gtx *layout.Context, theme *material.Theme, editor *widget.Editor) {
	layout.Flex{Axis: layout.Vertical}.Layout(*gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.UniformInset(unit.Dp(theme.TextSize)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return material.Editor(theme, editor, "Enter Text ...").Layout(gtx)
			})
		}),
	)

}
