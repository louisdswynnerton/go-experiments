package main

import (
	"math"
	"runtime"

	"github.com/PerformLine/go-stockutil/colorutil"
	"github.com/telroshan/go-sfml/v2/graphics"
	"github.com/telroshan/go-sfml/v2/window"
)

func getNullRenderState() graphics.SfRenderStates {
	return (graphics.SfRenderStates)(graphics.SwigcptrSfRenderStates(0))
}

func init() { runtime.LockOSThread() }

func main() {
	vm := window.NewSfVideoMode()
	defer window.DeleteSfVideoMode(vm)

	vm.SetWidth(800)
	vm.SetHeight(600)
	vm.SetBitsPerPixel(32)

	x := 1.0
	y := 1.0
	z := 1.0

	rho := 28.0
	sigma := 10.0
	beta := 8.0 / 3.0

	dx := 0.0
	dy := 0.0
	dz := 0.0

	xPath := []float32{}
	yPath := []float32{}
	zPath := []float32{}

	steps := 100000
	dt := 0.001

	currentstep := 2

	for i := 0; i < steps; i++ {

		dx = (sigma * (y - x)) * dt
		dy = (x*(rho-z) - y) * dt
		dz = (x*y - beta*z) * dt

		x += dx
		y += dy
		z += dz

		xPath = append(xPath, float32(x))
		yPath = append(yPath, float32(y))
		zPath = append(zPath, float32(z))

	}

	vertices := graphics.SfVertexArray_create()
	graphics.SfVertexArray_setPrimitiveType(vertices, graphics.SfPrimitiveType(graphics.SfLineStrip))

	vertex := graphics.NewSfVertex()
	posVector := graphics.NewSfVector2f()

	angle := 0.0

	var offset float32 = 0.0
	hue := 0.0

	color := graphics.NewSfColor()
	color.SetA(255)

	/* Create the main window */
	cs := window.NewSfContextSettings()
	defer window.DeleteSfContextSettings(cs)
	w := graphics.SfRenderWindow_create(vm, "SFML window", uint(window.SfResize|window.SfClose), cs)
	defer window.SfWindow_destroy(w)

	ev := window.NewSfEvent()
	defer window.DeleteSfEvent(ev)

	/* Start the game loop */
	for window.SfWindow_isOpen(w) > 0 {
		/* Process events */
		for window.SfWindow_pollEvent(w, ev) > 0 {
			/* Close window: exit */
			if ev.GetEvType() == window.SfEventType(window.SfEvtClosed) {
				return
			}

			if ev.GetEvType() == window.SfEventType(window.SfEvtKeyPressed) && ev.GetKey().GetCode() == window.SfKeyCode(window.SfKeyJ) {
				offset -= 1
			}

			if ev.GetEvType() == window.SfEventType(window.SfEvtKeyPressed) && ev.GetKey().GetCode() == window.SfKeyCode(window.SfKeyL) {
				offset += 1
			}
		}
		

		angle += 0.001

		vertices := graphics.SfVertexArray_create()
		graphics.SfVertexArray_setPrimitiveType(vertices, graphics.SfPrimitiveType(graphics.SfLineStrip))

		currentstep += 10

		for i := 0; i < currentstep; i++ {

			xDraw := ((xPath[i] - offset) * float32(math.Cos(angle))) + offset
			zDraw := ((zPath[i] - offset) * float32(math.Sin(angle))) + offset

			posVector.SetX((xDraw + zDraw) * 10)
			posVector.SetY((yPath[i] * 10) + 275)
			vertex.SetPosition(posVector)

			hue = (float64(i) / float64(steps)) * 360.0

			r, g, b := colorutil.HsvToRgb(hue, 1, 1)

			color.SetR(r)
			color.SetG(g)
			color.SetB(b)

			vertex.SetColor(color)
			graphics.SfVertexArray_append(vertices, vertex)

		}

		graphics.SfRenderWindow_clear(w, graphics.GetSfBlack())

		graphics.SfRenderWindow_drawVertexArray(w, vertices, getNullRenderState())

		graphics.SfRenderWindow_display(w)
	}
}
