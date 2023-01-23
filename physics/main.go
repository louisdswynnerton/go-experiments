package main

import (
	"math"
	"runtime"

	"github.com/telroshan/go-sfml/v2/graphics"
	"github.com/telroshan/go-sfml/v2/window"
)

func init() { runtime.LockOSThread() }

func getNullRenderState() graphics.SfRenderStates {
	return (graphics.SfRenderStates)(graphics.SwigcptrSfRenderStates(0))
}

type VerletObject struct {
	position_current  graphics.SfVector2f
	position_previous graphics.SfVector2f
	acceleration      graphics.SfVector2f
	radius            float32
}

func (o VerletObject) updatePosition(dt float32) {
	displacement := makeVector2(
		o.position_current.GetX()-o.position_previous.GetX(),
		o.position_current.GetY()-o.position_previous.GetY(),
	)

	o.position_previous.SetX(float32(o.position_current.GetX()))

	o.position_previous.SetY(float32(o.position_current.GetY()))

	o.position_current.SetX(float32(o.position_current.GetX() + displacement.GetX() + o.acceleration.GetX()*dt*dt))

	o.position_current.SetY(float32(o.position_current.GetY() + displacement.GetY() + o.acceleration.GetY()*dt*dt))

	o.acceleration.SetX(0.0)
	o.acceleration.SetY(0.0)

}

func (o VerletObject) accelerate(acc graphics.SfVector2f) {
	o.acceleration.SetX(o.acceleration.GetX() + acc.GetX())
	o.acceleration.SetY(o.acceleration.GetY() + acc.GetY())
}

type Solver struct {
	gravity        graphics.SfVector2f
	objectRegister []VerletObject
}

func (solver Solver) update(dt float32) {
	solver.applyGravity()
	solver.applyConstraints()
	solver.solveCollisions()
	solver.updatePositions(dt)

}

func (solver Solver) applyGravity() {
	for _, object := range (solver).objectRegister {
		object.accelerate((solver).gravity)
	}
}

func (solver Solver) updatePositions(dt float32) {
	for _, object := range (solver).objectRegister {
		object.updatePosition(dt)
	}
}

func (solver Solver) draw(w *graphics.Struct_SS_sfRenderWindow) {
	circle := graphics.SfCircleShape_create()
	for _, object := range solver.objectRegister {

		graphics.SfCircleShape_setRadius(circle, object.radius)
		graphics.SfCircleShape_setFillColor(circle, graphics.GetSfWhite())
		graphics.SfCircleShape_setPosition(circle, object.position_current)
		graphics.SfCircleShape_setOrigin(circle, makeVector2(object.radius, object.radius))
		graphics.SfRenderWindow_drawCircleShape((*w), circle, getNullRenderState())
		//fmt.Println(object.position_current.GetY())

	}
}

func (solver Solver) applyConstraints() {
	position := makeVector2(600.0, 400.0)
	radius := 300.0
	for _, object := range (solver).objectRegister {
		dist_x := position.GetX() - object.position_current.GetX()
		dist_y := position.GetY() - object.position_current.GetY()

		dist := math.Sqrt(float64(dist_x*dist_x + dist_y*dist_y))

		if dist > radius-float64(object.radius) {
			object.position_current.SetX(position.GetX() - (dist_x/float32(dist))*float32(radius-float64(object.radius)))
			object.position_current.SetY(position.GetY() - (dist_y/float32(dist))*float32(radius-float64(object.radius)))
			//fmt.Print(object.position_current.GetY())
			//fmt.Print(" ")
			//fmt.Println(position.GetY() - dist_y/float32(dist)*float32(radius-float64(object.radius)))

		}

	}

}

func (solver Solver) solveCollisions() {
	for i, object_1 := range (solver).objectRegister {
		for k, object_2 := range (solver).objectRegister {
			if i != k {

				axis_x := (object_1.position_current.GetX() - object_2.position_current.GetX())
				axis_y := (object_1.position_current.GetY() - object_2.position_current.GetY())
				dist2 := axis_x*axis_x + axis_y*axis_y
				min_dist := object_1.radius + object_2.radius

				if dist2 < min_dist*min_dist {
					dist := math.Sqrt(float64(dist2))
					n_x := (axis_x) / float32(dist)
					n_y := (axis_y) / float32(dist)
					delta := 0.5 * (min_dist - float32(dist))
					object_1.position_current.SetX(object_1.position_current.GetX() + n_x*delta)
					object_1.position_current.SetY(object_1.position_current.GetY() + n_y*delta)

					object_2.position_current.SetX(object_2.position_current.GetX() - n_x*delta)
					object_2.position_current.SetY(object_2.position_current.GetY() - n_y*delta)
				}
			}
		}
	}
}

func makeVector2(x float32, y float32) graphics.SfVector2f {
	v := graphics.NewSfVector2f()
	v.SetX(x)
	v.SetY(y)
	return v
}

func makeVerletObject(position graphics.SfVector2f, radius float32) VerletObject {
	newObject := VerletObject{
		position_current:  makeVector2(position.GetX(), position.GetY()),
		position_previous: makeVector2(position.GetX(), position.GetY()),
		acceleration:      makeVector2(0.0, 0.0),
		radius:            radius,
	}
	return newObject
}

func main() {

	gravity := graphics.NewSfVector2f()
	gravity.SetX(0.0)
	gravity.SetY(1000.0)

	solver := Solver{
		gravity,
		[]VerletObject{},
	}

	center := graphics.NewSfVector2f()

	center.SetX(600.0)
	center.SetY(400.0)

	zeroVector := graphics.NewSfVector2f()
	zeroVector.SetX(0.0)
	zeroVector.SetY(0.0)

	vm := window.NewSfVideoMode()
	defer window.DeleteSfVideoMode(vm)
	vm.SetWidth(1200)
	vm.SetHeight(800)
	vm.SetBitsPerPixel(32)
	cs := window.NewSfContextSettings()
	defer window.DeleteSfContextSettings(cs)
	w := graphics.SfRenderWindow_create(vm, "Louis is great", uint(window.SfTitlebar|window.SfClose), cs)
	defer window.SfWindow_destroy(w)

	ev := window.NewSfEvent()
	defer window.DeleteSfEvent(ev)

	//window.SfWindow_setFramerateLimit(w, 60)

	testcircle := graphics.SfCircleShape_create()
	graphics.SfCircleShape_setPointCount(testcircle, 100)
	graphics.SfCircleShape_setFillColor(testcircle, graphics.GetSfBlack())
	graphics.SfCircleShape_setRadius(testcircle, 300.0)
	graphics.SfCircleShape_setPosition(testcircle, makeVector2(600.0, 400.0))
	graphics.SfCircleShape_setOrigin(testcircle, makeVector2(300.0, 300.0))

	for window.SfWindow_isOpen(w) > 0 {
		/* Process events */
		for window.SfWindow_pollEvent(w, ev) > 0 {
			/* Close window: exit */
			if ev.GetEvType() == window.SfEventType(window.SfEvtClosed) {
				return
			}

			if ev.GetEvType() == window.SfEventType(window.SfEvtKeyPressed) && ev.GetKey().GetCode() == window.SfKeyCode(window.SfKeySpace) {

				solver.objectRegister = append(solver.objectRegister, makeVerletObject(makeVector2(500.0, 400.0), 15.0))
			}

		}

		solver.update(0.01)

		graphics.SfRenderWindow_clear(w, graphics.GetSfCyan())
		graphics.SfRenderWindow_drawCircleShape((w), testcircle, getNullRenderState())

		solver.draw(&w)

		graphics.SfRenderWindow_display(w)

	}

}
