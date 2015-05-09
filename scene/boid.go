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
	BOID_DEFAULT_SIZE int32 = 10
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
	b.Velocity.X = (rand.Float64() * 2.0) - 1.0
	b.Velocity.Y = (rand.Float64() * 2.0) - 1.0
	b.PixelSize = BOID_DEFAULT_SIZE

	return b
}

func (b *Boid) drawBoid(w, h int32, renderer *sdl.Renderer) {
	x, y := b.boidPosWithinBounds(w, h)

	renderer.SetDrawColor(255, 255, 255, 255)

	// direction
	// a := b.Velocity.HeadingAngle()
	// log.Println(b.Velocity.X, b.Velocity.Y, a)

	// TODO: use heading angle
	renderer.DrawLine(int(x), int(y), int(x)+3, int(y)+3)
}

func (b Boid) boidPosWithinBounds(w, h int32) (x, y int32) {
	// adjust the coordinates to reflect the centre of the boid
	xf := float64(b.X) * float64(w-(b.PixelSize/2))
	yf := float64(b.Y) * float64(h-(b.PixelSize/2))

	return int32(xf), int32(yf)
}
