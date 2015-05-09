package scene

import "math"

type Vector struct {
	X float64
	Y float64
}

func (v Vector) Magnitude() float64 {
	return math.Sqrt((v.X * v.X) + (v.Y * v.Y))
}

func (v *Vector) Add(b Vector) {
	v.X = v.X + b.X
	v.Y = v.Y + b.Y
	return
}

func (v *Vector) Sub(b Vector) {
	v.X = v.X - b.X
	v.Y = v.Y - b.Y
	return
}

func (v *Vector) Div(a float64) {
	v.X /= a
	v.Y /= a
	return
}

func VectorAdd(a Vector, b Vector) Vector {
	return Vector{a.X + b.X, a.Y + b.Y}
}

func VectorSub(a Vector, b Vector) Vector {
	return Vector{a.X - b.X, a.Y - b.Y}
}

func VectorDiv(v Vector, a float64) Vector {
	return Vector{v.X / a, v.Y / a}
}

func VectorLimitSpeed(v Vector, max_speed float64) Vector {

	mod := v.Magnitude()
	x := v.X
	y := v.Y
	if mod > max_speed {
		x = (v.X / mod) * max_speed
		y = (v.Y / mod) * max_speed

	}

	return Vector{x, y}
}

func (v *Vector) LimitSpeed(max_speed float64) {

	mod := v.Magnitude()
	if mod > max_speed {
		v.X = (v.X / mod) * max_speed
		v.Y = (v.Y / mod) * max_speed

	}
}

func (v *Vector) Normalise() {
	m := v.Magnitude()
	v.X = v.X / m
	v.Y = v.Y / m
}

func (v Vector) HeadingAngle() float64 {
	return math.Atan2(v.Y, v.X)
}
