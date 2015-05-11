package scene

import (
	"fmt"
	"log"
	"time"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/sdl_ttf"
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

	DEFAULT_LEADER_COUNT = 3
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

	s.ShowActivePads = false
	s.ShowGrid = false
	s.ShowHUD = false

	s.Triggers = NewTriggerGrid(8, 0.8)

	return s
}

func (s *Scene) AddLeader() {

	if s.LeaderCount < len(s.Boids) {

		s.Boids[s.LeaderCount].BoidKind = BoidKind_LEADER
		s.LeaderCount += 1
	}
}

func (s *Scene) RemoveLeader() {
	if s.LeaderCount > 0 {

		s.Boids[s.LeaderCount].BoidKind = BoidKind_NORMAL
		s.LeaderCount -= 1
	} else {
		s.Boids[s.LeaderCount].BoidKind = BoidKind_NORMAL
		s.LeaderCount = 0
	}
}

func (s *Scene) AllLeaders() {

	for j := s.LeaderCount; j < len(s.Boids); j++ {

		s.Boids[j].BoidKind = BoidKind_LEADER
		s.LeaderCount = len(s.Boids) - 1

	}

}

func (s *Scene) RestoreDefault() {
	s.Speed = DEFAULT_SPEED
	s.Distance = DEFAULT_DISTANCE
	s.Radius = DEFAULT_RADIUS
	s.CohesionWeight = DEFAULT_COHESION_WEIGHT
	s.AlignmentWeight = DEFAULT_ALIGNMENT_WEIGHT
	s.SeparationWeight = DEFAULT_SEPARATION_WEIGHT

	// restore leader count
	for j := 0; j <= DEFAULT_LEADER_COUNT; j++ {
		s.Boids[j].BoidKind = BoidKind_LEADER

	}

	for j := DEFAULT_LEADER_COUNT; j < len(s.Boids); j++ {
		s.Boids[j].BoidKind = BoidKind_NORMAL

	}

	s.LeaderCount = DEFAULT_LEADER_COUNT
}

func (s *Scene) Draw(w, h int32, renderer *sdl.Renderer, font *ttf.Font) {
	s.drawTriggers(w, h, renderer)
	s.drawBoids(w, h, renderer)
	s.drawHUD(w, h, renderer, font)
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

		// update boid rules
		s.Boids[i].Velocity.X = s.Boids[i].Velocity.X + cohesion.X + separation.X + alignment.X
		s.Boids[i].Velocity.Y = s.Boids[i].Velocity.Y + cohesion.Y + separation.Y + alignment.Y

		// limit spped
		s.Boids[i].Velocity.LimitSpeed(s.Speed)

		// update positions according to new velocity
		s.Boids[i].X += b.Velocity.X
		s.Boids[i].Y += b.Velocity.Y

		// wrap around borders

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
			// calculate distance between boids
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
		// average Sum
		sum.X /= float64(count)
		sum.Y /= float64(count)

		// vector pointing from the location to the desired target
		desired := Vector{}
		desired.X = sum.X - s.Boids[i].X
		desired.Y = sum.Y - s.Boids[i].Y

		// distance from the target is the magnitude  of the vector
		d := desired.Magnitude()

		if d > 0 {
			desired.Normalise()
			desired.X = desired.X * s.Speed
			desired.Y = desired.Y * s.Speed

			// and Steer
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
			// calculate distance between boids

			distance := Vector{}
			distance.X = s.Boids[j].X - s.Boids[i].X
			distance.Y = s.Boids[j].Y - s.Boids[i].Y

			d := distance.Magnitude()

			if (d < s.Distance) && (d > 0.0) {

				diff := Vector{}
				diff.X = s.Boids[i].X - s.Boids[j].X
				diff.Y = s.Boids[i].Y - s.Boids[j].Y

				diff.Normalise()

				// weight by distance
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
			// calculate distance between boids
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

func (s *Scene) UpdateTriggers() []TriggerEvent {

	events := make([]TriggerEvent, 0)
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

	// change in state should result in on of off messages
	for j := range s.Triggers {

		te := TriggerEvent{}
		te.TriggerEventType = s.Triggers[j].StateTransition()
		te.Note = s.Triggers[j].Note
		te.Channel = s.Triggers[j].Channel

		switch te.TriggerEventType {
		case TriggerEventType_OFF:
			events = append(events, te)
		case TriggerEventType_ON:
			events = append(events, te)
		}
	}

	return events
}

func (s *Scene) AddBoid(b Boid) {
	s.Boids = append(s.Boids, b)
}

func (s *Scene) String() string {
	return fmt.Sprintf("Boid count: %d", len(s.Boids))
}

func (s *Scene) drawHUD(w, h int32, renderer *sdl.Renderer, font *ttf.Font) {

	if !s.ShowHUD {
		return
	}
	colorRedSemiTrans := sdl.Color{180, 0, 0, 100}
	colorRed := sdl.Color{180, 0, 0, 255}

	var lineHeight int32 = 24 + 8
	var xOffset int32 = 20
	var yOffset int32 = 20

	renderText(fmt.Sprintf("Boids:      %d", len(s.Boids)), xOffset, yOffset, &colorRedSemiTrans, renderer, font)
	renderText(fmt.Sprintf("Speed:      %.3f", s.Speed), xOffset, yOffset+(lineHeight*2), &colorRed, renderer, font)
	renderText(fmt.Sprintf("Distance:   %.3f", s.Distance), xOffset, yOffset+(lineHeight*3), &colorRed, renderer, font)
	renderText(fmt.Sprintf("Radius:     %.3f", s.Radius), xOffset, yOffset+(lineHeight*4), &colorRed, renderer, font)
	renderText(fmt.Sprintf("Cohesion:   %.3f", s.CohesionWeight), xOffset, yOffset+(lineHeight*5), &colorRed, renderer, font)
	renderText(fmt.Sprintf("Alignment:  %.3f", s.AlignmentWeight), xOffset, yOffset+(lineHeight*6), &colorRed, renderer, font)
	renderText(fmt.Sprintf("Separation: %.3f", s.SeparationWeight), xOffset, yOffset+(lineHeight*7), &colorRed, renderer, font)
}

func renderText(text string, x, y int32, color *sdl.Color, renderer *sdl.Renderer, font *ttf.Font) {
	surface := font.RenderText_Solid(text, *color)
	texture, err := renderer.CreateTextureFromSurface(surface)
	if err != nil {
		log.Fatalf("Failed to create texture: %s\n", err)
	}
	src := sdl.Rect{0, 0, surface.W, surface.H}
	dst := sdl.Rect{x, y, surface.W, surface.H}

	renderer.Copy(texture, &src, &dst)
}
