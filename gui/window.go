package gui

import (
	"fmt"
	"image"
	"image/color"
	"os"

	"gioui.org/app"
	"gioui.org/io/key"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
)

const title = "Generated diagram preview..."

type GUI struct {
	image paint.ImageOp
	title string
}

func NewGUI(img image.Image) *GUI {
	return &GUI{
		image: paint.NewImageOp(img),
		title: title,
	}
}

func (ui *GUI) Draw() error {
	w := new(app.Window)

	w.Option(app.Size(
		unit.Dp(ui.image.Size().X),
		unit.Dp(ui.image.Size().Y),
	), app.Title(ui.title))

	if err := ui.run(w); err != nil {
		defer func() {
			os.Exit(0)
		}()
		return fmt.Errorf("GUI rendering error: %w", err)
	}

	return nil
}

func (ui *GUI) run(w *app.Window) error {
	var ops op.Ops

	for {
		switch e := w.Event().(type) {
		case app.FrameEvent:
			gtx := app.NewContext(&ops, e)
			for {
				event, ok := gtx.Event(key.Filter{
					Name: key.NameEscape,
				})
				if !ok {
					break
				}
				switch event := event.(type) {
				case key.Event:
					switch event.Name {
					case key.NameEscape:
						w.Perform(system.ActionClose)
					}
				}
			}

			ui.drawDiagram(gtx)
			e.Frame(gtx.Ops)
		case app.DestroyEvent:
			return e.Err
		}
	}
}

func (ui *GUI) drawDiagram(gtx layout.Context) {
	layout.Stack{}.Layout(gtx,
		layout.Stacked(func(gtx layout.Context) layout.Dimensions {
			paint.FillShape(gtx.Ops, color.NRGBA{R: 0xFF, G: 0xFF, B: 0xFF, A: 0xff},
				clip.Rect{Max: gtx.Constraints.Max}.Op(),
			)

			return layout.UniformInset(unit.Dp(0)).Layout(gtx,
				func(gtx layout.Context) layout.Dimensions {
					widget.Image{
						Src:   ui.image,
						Scale: 1 / float32(unit.Dp(1)),
						Fit:   widget.Contain,
					}.Layout(gtx)

					return layout.Dimensions{Size: gtx.Constraints.Max}
				})
		}),
	)
}
