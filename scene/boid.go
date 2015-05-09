package scene

import (
	"math/rand"

	"github.com/veandco/go-sdl2/sdl"
)

type BoidKind int32

const (
	BoidKind_NORMAL   BoidKind = 0
	BoidKind_LEADER   BoidKind = 1
	BoidKind_MUSHROOM BoidKind = 2
)

const (
	BOID_DEFAULT_SIZE int32 = 3
)

type Boid struct {
	X, Y      float64
	Velocity  Vector
	BoidKind  BoidKind
	PixelSize int32
}

func NewRandomBoid() Boid {
	b := Boid{}

	b.X = rand.Float64()
	b.Y = rand.Float64()
	b.Velocity = Vector{}
	b.Velocity.X = rand.Float64()
	b.Velocity.Y = rand.Float64()
	b.PixelSize = BOID_DEFAULT_SIZE

	return b
}

func (b *Boid) drawBoid(w, h int32, renderer *sdl.Renderer) {
	x, y := b.boidPosWithinBounds(w, h)

	renderer.SetDrawColor(255, 255, 255, 255)

	r := &sdl.Rect{}
	r.H = b.PixelSize
	r.W = b.PixelSize
	r.X = x
	r.Y = y

	renderer.FillRect(r)
}

func (b Boid) boidPosWithinBounds(w, h int32) (x, y int32) {
	// adjust the coordinates to reflect the centre of the boid
	xf := float64(b.X) * float64(w-(b.PixelSize/2))
	yf := float64(b.Y) * float64(h-(b.PixelSize/2))

	return int32(xf), int32(yf)
}
