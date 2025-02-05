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
		window.Option(app.Size(unit.Dp(800), unit.Dp(600)))

		err := run(window)
		if err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()
	app.Main()

}

func run(window *app.Window) error {
	theme := material.NewTheme() //creates a new theme for styling the ui
	var ops op.Ops
	editor := new(widget.Editor)
	editor.SingleLine = false
	var list widget.List        //creates a variable which is used to manage a scrollable list of items (for scrollbar)
	list.Axis = layout.Vertical //makes it so the list scrolls vertically
	list.ScrollToEnd = false
	for {
		switch e := window.Event().(type) {
		case app.DestroyEvent:
			return e.Err
		case app.FrameEvent:
			gtx := app.NewContext(&ops, e)

			layout.Flex{
				Axis: layout.Vertical,
			}.Layout(gtx,
				layout.Rigid(func(gtx layout.Context) layout.Dimensions {
					return addEditor(&list, gtx, theme, editor)
				}),
			)
			e.Frame(gtx.Ops)
		}
	}
}

func addEditor(list *widget.List, gtx layout.Context, theme *material.Theme, editor *widget.Editor) layout.Dimensions { //renders the text editor widget inside a list

	return material.List(theme, list).Layout(gtx, 1, func(gtx layout.Context, index int) layout.Dimensions { //renders a scrollable list
		margins := layout.Inset{
			Top:    unit.Dp(8),
			Bottom: unit.Dp(8),
			Left:   unit.Dp(8),
			Right:  unit.Dp(8),
		}
		return margins.Layout(gtx, func(gtx layout.Context) layout.Dimensions { //applies the previous margins
			gtx.Constraints.Min.Y = gtx.Dp(unit.Dp(300))
			return material.Editor(theme, editor, "Enter Text ...").Layout(gtx)
		})
	})
}
