package scene

import (
	"fmt"
	"time"

	"github.com/veandco/go-sdl2/sdl"
)

type Scene struct {
	Boids         []Boid
	Width, Height int32
	TheBigBang    time.Time
}

type Trigger struct {
}

func NewScene(boidCount int32, w, h int32) Scene {
	s := Scene{}
	s.TheBigBang = time.Now()
	s.Width = w
	s.Height = h
	s.Boids = make([]Boid, boidCount)
	for i := range s.Boids {
		s.Boids[i] = NewRandomBoid()
	}

	return s
}

func (s *Scene) Draw(w, h int32, renderer *sdl.Renderer) {
	s.drawBoids(w, h, renderer)
}

func (s *Scene) drawBoids(w, h int32, renderer *sdl.Renderer) {
	for i := range s.Boids {
		s.Boids[i].drawBoid(w, h, renderer)
	}
}

func (s *Scene) UpdateBoids() {
	// TODO: update Boids
}

func (s *Scene) UpdateTriggers() {
	// TODO: update triggers and Trigger state
}

func (s *Scene) AddBoid(b Boid) {
	s.Boids = append(s.Boids, b)
}

func (s *Scene) String() string {
	return fmt.Sprintf("Boid count: %d", len(s.Boids))
}
