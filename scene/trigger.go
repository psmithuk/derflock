package scene

import "github.com/veandco/go-sdl2/sdl"

type Trigger struct {
	Channel        int32
	Note           int32
	X1, Y1, X2, Y2 float64
	Active         bool
	PreviousState  bool
}

type TriggerEventType int32

type TriggerEvent struct {
	TriggerEventType TriggerEventType
	Channel          int32
	Note             int32
}

const (
	TriggerEventType_NONE TriggerEventType = 0
	TriggerEventType_ON   TriggerEventType = 1
	TriggerEventType_OFF  TriggerEventType = 2
)

func NewTriggerGrid(squareSize int, density float64) []Trigger {
	t := make([]Trigger, squareSize*squareSize)

	// the space which the trigger will occupy
	triggerVoid := 1.0 / float64(squareSize)
	offset := 0.005

	for row := 0; row < squareSize; row++ {
		for col := 0; col < squareSize; col++ {
			trigger := Trigger{}
			trigger.X1 = (float64(row) * triggerVoid) + offset
			trigger.Y1 = (float64(col) * triggerVoid) + offset
			trigger.X2 = (trigger.X1 + (density * triggerVoid)) + offset
			trigger.Y2 = (trigger.Y1 + (density * triggerVoid)) + offset
			trigger.Note = int32((row * 16) + col)
			trigger.Channel = 1
			t[(row*squareSize)+col] = trigger
		}
	}

	return t
}

func (t TriggerEventType) String() string {
	switch t {
	case TriggerEventType_ON:
		return "ON"
	case TriggerEventType_OFF:
		return "OFF"
	default:
		return "-"
	}
}

func (t *Trigger) StateTransition() TriggerEventType {
	e := TriggerEventType_NONE
	if t.Active != t.PreviousState {
		if t.Active {
			e = TriggerEventType_ON
		} else {
			e = TriggerEventType_OFF
		}
		t.PreviousState = t.Active
	}
	return e
}

func (t *Trigger) drawTrigger(w, h int32, renderer *sdl.Renderer, showActive, showGrid bool) {
	r := &sdl.Rect{
		X: int32(float64(w) * t.X1),
		Y: int32(float64(h) * t.Y1),
		W: int32((float64(w) * t.X2) - (float64(w) * t.X1)),
		H: int32((float64(h) * t.Y2) - (float64(h) * t.Y1)),
	}

	if t.Active && showActive {
		renderer.SetDrawColor(240, 240, 100, 255)
		renderer.FillRect(r)
	}

	if showGrid {
		renderer.SetDrawColor(100, 100, 100, 255)
		renderer.DrawRect(r)

	}
}
