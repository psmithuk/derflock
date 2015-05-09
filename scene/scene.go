package scene

import (
	"fmt"
	"time"

	"github.com/veandco/go-sdl2/sdl"
)

type Scene struct {
	Boids    []Boid
	Speed    float64
	Distance float64
	Radius   float64

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

	s.Speed = 0.005
	s.Distance = 0.02
	s.Radius = 0.1

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
	cohesion := Vector{}
	//alignment := Vector{}
	//separation := Vector{}

	//For each boid update positions
	for i, b := range s.Boids {

		cohesion = s.CohesionForBoid(i)

		//Update boid rules
		//First cohesion
		s.Boids[i].Velocity.Add(cohesion)

		//limit spped
		//s.Boids[i].Velocity.LimitSpeed(s.Speed)

		//Update position
		s.Boids[i].X += (b.Velocity.X * s.Speed)
		s.Boids[i].Y += (b.Velocity.Y * s.Speed)

		//Limit borders

		if s.Boids[i].X >= 1.0 {
			s.Boids[i].X = 0

		} else if s.Boids[i].X <= 0.0 {
			s.Boids[i].X = 1
		}

		if s.Boids[i].Y >= 1.0 {
			s.Boids[i].Y = 0

		} else if s.Boids[i].Y <= 0.0 {
			s.Boids[i].Y = 1
		}

	}

}

func (s *Scene) CohesionForBoid(i int) Vector {
	sum := Vector{}
	var count = 0

	for j := 0; j < len(s.Boids); j++ {

		if j != i {
			//calculate distance between boids
			distance := Vector{}
			distance.X = s.Boids[j].X - s.Boids[i].X
			distance.Y = s.Boids[j].Y - s.Boids[i].Y
			d := distance.Magnitude()

			if d > 0 && d < s.Radius {

				sum.X += s.Boids[j].X
				sum.Y += s.Boids[j].Y
				count += 1
			}
		}
	}

	if count > 0 {
		//Average Sum
		sum.X /= float64(count)
		sum.Y /= float64(count)

		//And Steer
		steer := Vector{}

		//a vector pointing from the location to the target
		desired := Vector{sum.X, sum.Y}
		location := Vector{s.Boids[i].X, s.Boids[i].Y}

		desired.Sub(location)
		d := desired.Magnitude()
		if d > 0 {

			desired.Normalise()
			desired.X = desired.X * s.Speed
			desired.Y = desired.Y * s.Speed

			steer = desired
			steer.Sub(s.Boids[i].Velocity)
			steer.LimitSpeed(s.Speed)

			return steer
		}

	}

	return sum

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
