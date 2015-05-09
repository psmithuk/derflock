package scene

import (
	"fmt"
	"time"

	"github.com/veandco/go-sdl2/sdl"
)

type Scene struct {
	Boids            []Boid
	Triggers         []Trigger
	Speed            float64
	Distance         float64
	Radius           float64
	CohesionWeight   float64
	AlignmentWeight  float64
	SeparationWeight float64
	LeaderCount      int

	ShowGrid       bool
	ShowHUD        bool
	ShowActivePads bool

	Width, Height int32
	TheBigBang    time.Time
}

const (
	DEFAULT_SPEED    = 0.004
	DEFAULT_DISTANCE = 0.01
	DEFAULT_RADIUS   = 0.1

	DEFAULT_COHESION_WEIGHT   = 0.1
	DEFAULT_ALIGNMENT_WEIGHT  = 1
	DEFAULT_SEPARATION_WEIGHT = 0.1
)

func NewScene(boidCount int32, w, h int32) Scene {
	s := Scene{}
	s.TheBigBang = time.Now()
	s.Width = w
	s.Height = h
	s.Boids = make([]Boid, boidCount)
	for i := range s.Boids {
		s.Boids[i] = NewRandomBoid()
	}

	s.Speed = DEFAULT_SPEED
	s.Distance = DEFAULT_DISTANCE
	s.Radius = DEFAULT_RADIUS

	s.CohesionWeight = DEFAULT_COHESION_WEIGHT
	s.AlignmentWeight = DEFAULT_ALIGNMENT_WEIGHT
	s.SeparationWeight = DEFAULT_SEPARATION_WEIGHT

	s.ShowActivePads = true
	s.ShowGrid = true
	s.ShowHUD = true

	s.Triggers = NewTriggerGrid(8, 0.8)

	return s
}
func (s *Scene) AddLeader() {

	s.Boids[s.LeaderCount].BoidKind = BoidKind_LEADER
	s.LeaderCount += 1
}

func (s *Scene) AllLeaders() {

	for j := s.LeaderCount; j < len(s.Boids); j++ {

		s.Boids[j].BoidKind = BoidKind_LEADER
		s.LeaderCount += 1

	}

}

func (s *Scene) RestoreDefault() {
	s.Speed = DEFAULT_SPEED
	s.Distance = DEFAULT_DISTANCE
	s.Radius = DEFAULT_RADIUS
	s.CohesionWeight = DEFAULT_COHESION_WEIGHT
	s.AlignmentWeight = DEFAULT_ALIGNMENT_WEIGHT
	s.SeparationWeight = DEFAULT_SEPARATION_WEIGHT
}

func (s *Scene) Draw(w, h int32, renderer *sdl.Renderer) {
	s.drawTriggers(w, h, renderer)
	s.drawBoids(w, h, renderer)
}

func (s *Scene) drawBoids(w, h int32, renderer *sdl.Renderer) {
	for i := range s.Boids {
		s.Boids[i].drawBoid(w, h, renderer)
	}
}

func (s *Scene) drawTriggers(w, h int32, renderer *sdl.Renderer) {
	for i := range s.Triggers {
		s.Triggers[i].drawTrigger(w, h, renderer, s.ShowActivePads, s.ShowGrid)
	}
}

func (s *Scene) UpdateBoids() {

	//For each boid update positions
	for i, b := range s.Boids {
		cohesion := Vector{}
		separation := Vector{}
		alignment := Vector{}

		cohesion = s.CohesionForBoid(i)
		separation = s.SeparationForBoid(i)
		alignment = s.AlignmentForBoid(i)

		//Update boid rules
		//First cohesion
		s.Boids[i].Velocity.X = s.Boids[i].Velocity.X + cohesion.X + separation.X + alignment.X
		s.Boids[i].Velocity.Y = s.Boids[i].Velocity.Y + cohesion.Y + separation.Y + alignment.Y
		//limit spped
		s.Boids[i].Velocity.LimitSpeed(s.Speed)

		//Update position
		s.Boids[i].X += b.Velocity.X
		s.Boids[i].Y += b.Velocity.Y

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

		//a vector pointing from the location to the desired target
		desired := Vector{}
		desired.X = sum.X - s.Boids[i].X
		desired.Y = sum.Y - s.Boids[i].Y

		//Distance from the target is the magnitude  of the vector
		d := desired.Magnitude()
		if d > 0 {

			desired.Normalise()
			desired.X = desired.X * s.Speed
			desired.Y = desired.Y * s.Speed

			//And Steer
			steer := Vector{}

			steer.X = (desired.X - s.Boids[i].Velocity.X) * s.CohesionWeight
			steer.Y = (desired.Y - s.Boids[i].Velocity.Y) * s.CohesionWeight

			return steer
		} else {
			return Vector{0.0, 0.0}
		}

	}

	sum.X *= s.CohesionWeight
	sum.Y *= s.CohesionWeight
	return sum

}

func (s *Scene) SeparationForBoid(i int) Vector {

	sum := Vector{}
	var count = 0

	for j := 0; j < len(s.Boids); j++ {

		if j != i {
			//calculate distance between boids
			distance := Vector{}
			distance.X = s.Boids[j].X - s.Boids[i].X
			distance.Y = s.Boids[j].Y - s.Boids[i].Y
			d := distance.Magnitude()

			if (d < s.Distance) && (d > 0.0) {

				diff := Vector{}
				diff.X = s.Boids[i].X - s.Boids[j].X
				diff.Y = s.Boids[i].Y - s.Boids[j].Y

				diff.Normalise()

				//weight by distance
				diff.X /= d
				diff.Y /= d

				sum.X += diff.X
				sum.Y += diff.Y

				count += 1
			}
		}
	}

	if count > 0 {
		sum.X /= float64(count)
		sum.Y /= float64(count)

	}

	sum.X *= s.SeparationWeight
	sum.Y *= s.SeparationWeight
	return sum

}

func (s *Scene) AlignmentForBoid(i int) Vector {

	sum := Vector{}
	var count = 0

	for j := 0; j < len(s.Boids); j++ {

		if j != i {
			//calculate distance between boids
			distance := Vector{}
			distance.X = s.Boids[j].X - s.Boids[i].X
			distance.Y = s.Boids[j].Y - s.Boids[i].Y
			d := distance.Magnitude()

			if (d < s.Radius) && (d > 0.0) {

				sum.X += s.Boids[j].Velocity.X
				sum.Y += s.Boids[j].Velocity.Y
				count += 1
			}
		}
	}

	if count > 0 {
		sum.X /= float64(count)
		sum.Y /= float64(count)

	}

	sum.X *= s.AlignmentWeight
	sum.Y *= s.AlignmentWeight
	return sum

}

func (s *Scene) UpdateTriggers() {

	// clear trigger state
	for j := range s.Triggers {
		s.Triggers[j].Active = false
	}

	for i, boid := range s.Boids {
		// only trigger the leaders
		if s.Boids[i].BoidKind == BoidKind_LEADER {
			for j, t := range s.Triggers {
				if boid.X >= t.X1 && boid.X <= t.X2 && boid.Y >= t.Y1 && boid.Y <= t.Y2 {
					s.Triggers[j].Active = true
				}
			}
		}
	}
}

func (s *Scene) AddBoid(b Boid) {
	s.Boids = append(s.Boids, b)
}

func (s *Scene) String() string {
	return fmt.Sprintf("Boid count: %d", len(s.Boids))
}
