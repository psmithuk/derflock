package scene

import "math"

type Vector struct {
	X float64
	Y float64
}

func (v Vector) Magnitude() float64 {
	return math.Sqrt((v.X * v.X) + (v.Y * v.Y))
}

func (v Vector) Add(b Vector) {
	v.X = v.X + b.X
	v.Y = v.Y + b.Y
	return
}

func (v Vector) Sub(b Vector) {
	v.X = v.X - b.X
	v.Y = v.Y - b.Y
	return
}

func (v Vector) Div(a float64) {
	v.X /= a
	v.Y /= a
	return
}

func (v Vector) LimitSpeed(max_speed float64) {

	mod := v.Magnitude()
	if mod > max_speed {
		v.X = (v.X / mod) * max_speed
		v.Y = (v.Y / mod) * max_speed

	}
	return
}

func (v Vector) Normalise() Vector {
	m := v.Magnitude()
	return Vector{
		X: v.X / m,
		Y: v.Y / m,
	}
}
