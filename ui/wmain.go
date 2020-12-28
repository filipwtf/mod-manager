package ui

import (
	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget/material"
	"github.com/filipwtf/filips-installer/config"
	"image"
	"image/color"
)

func (main *mainWidget) managerLayout(th *material.Theme, gtx ctx, mods []Mod, full bool) dim {
	var size float32
	if full {
		size = 0.68
	} else {
		size = 0.92
	}

	height := float32(gtx.Constraints.Max.Y) * size
	gtx.Constraints.Max.Y = int(height)
	widgets := []layout.Widget{
		material.H4(th, "Mod Manager").Layout,
		drawSplitter(),
		material.Body1(th, "Mod Name").Layout,
		// TODO Table title
		drawSplitter(),
	}

	for _, mod := range mods {
		mod := mod
		widgets = append(widgets, mod.modWidget(th))
	}

	return list.Layout(gtx, len(widgets), func(gtx ctx, i int) dim {
		return layout.Inset{
			Top: unit.Dp(0),
		}.Layout(gtx, widgets[i])
	})
}

func (main *mainWidget) installerLayout(th *material.Theme, gtx ctx, full bool) dim {
	var size float32
	if full {
		size = 0.68
	} else {
		size = 0.92
	}
	height := float32(gtx.Constraints.Max.Y) * size
	gtx.Constraints.Max.Y = int(height)
	widgets := []layout.Widget{
		material.H4(th, "Mod Installer").Layout,
	}

	return list.Layout(gtx, len(widgets), func(gtx ctx, i int) dim {
		return layout.Inset{
			Top:  unit.Dp(0),
			Left: unit.Dp(10),
		}.Layout(gtx, widgets[i])
	})
}

func drawSplitter() layout.Widget {
	return func(gtx ctx) dim {
		ops := gtx.Ops
		clip.Rect{Max: image.Pt(int(config.Width.V), 100)}.Add(ops)
		paint.ColorOp{Color: color.NRGBA{R: 0x00, A: 0xFF}}.Add(ops)
		paint.PaintOp{}.Add(ops)
		return layoutWidget(1, 1)
	}
}
