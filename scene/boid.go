package scene

import (
	"math"
	"math/rand"
	"time"

	"github.com/veandco/go-sdl2/sdl"
)

type BoidKind int32

const (
	BoidKind_NORMAL   BoidKind = 0
	BoidKind_LEADER   BoidKind = 1
	BoidKind_MUSHROOM BoidKind = 2
)

const (
	BOID_DEFAULT_SIZE int32 = 5
)

type Boid struct {
	X, Y      float64
	Velocity  Vector
	BoidKind  BoidKind
	PixelSize int32
}

func NewRandomBoid() Boid {
	b := Boid{}

	rand.Seed(time.Now().UnixNano())

	b.X = rand.NormFloat64()
	b.Y = rand.NormFloat64()
	b.Velocity = Vector{}
	b.Velocity.X = rand.NormFloat64()
	b.Velocity.Y = rand.NormFloat64()
	b.PixelSize = BOID_DEFAULT_SIZE

	return b
}

func (b *Boid) drawBoid(w, h int32, renderer *sdl.Renderer) {
	x, y := b.boidPosWithinBounds(w, h)

	switch b.BoidKind {
	case BoidKind_LEADER:
		renderer.SetDrawColor(255, 255, 0, 255)
	default:
		renderer.SetDrawColor(255, 255, 255, 255)

	}

	// direction
	a := b.Velocity.HeadingAngle()

	direction := Vector{
		X: math.Cos(a),
		Y: math.Sin(a),
	}

	renderer.DrawLine(int(x),
		int(y),
		int(float64(x)+float64(b.PixelSize)*direction.X),
		int(float64(y)+float64(b.PixelSize)*direction.Y))
}

func (b Boid) boidPosWithinBounds(w, h int32) (x, y int32) {
	// adjust the coordinates to reflect the centre of the boid
	xf := float64(b.X) * float64(w-(b.PixelSize/2))
	yf := float64(b.Y) * float64(h-(b.PixelSize/2))

	return int32(xf), int32(yf)
}
