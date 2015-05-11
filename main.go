package main

import (
	"log"
	"math/rand"
	"runtime"
	"time"

	"github.com/psmithuk/derflock/scene"
	"github.com/rakyll/portmidi"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/sdl_ttf"
)

var winTitle string = "Der Flock"
var width, height int32 = 720, 720
var font *ttf.Font

func init() {
	runtime.GOMAXPROCS(3)
	runtime.LockOSThread()

	rand.Seed(time.Now().UnixNano())

	portmidi.Initialize()
}

func main() {
	var window *sdl.Window
	var renderer *sdl.Renderer
	var event sdl.Event
	var err error

	var deviceId portmidi.DeviceId
	var outputStream *portmidi.Stream

	useMidi := false

	midiCount := portmidi.CountDevices()
	for i := 0; i < midiCount; i++ {
		d := portmidi.GetDeviceInfo(portmidi.DeviceId(i))
		log.Println(i, d.Name, d.IsOutputAvailable)
		if d.Name == "IAC Driver DerFlock" && d.IsOutputAvailable == true {
			deviceId = portmidi.DeviceId(i)
			log.Printf("%v\n", deviceId)
			useMidi = true
		}
	}

	if useMidi {
		outputStream, err = portmidi.NewOutputStream(deviceId, 1024, 0)
		if err != nil {
			log.Fatalf("Failed open new midi output: %s", err)
		}
	}

	running := true

	err = sdl.Init(sdl.INIT_VIDEO)
	if err != nil {
		log.Fatalf("Failed to init: %s\n", err)
	}

	window, err = sdl.CreateWindow(winTitle, sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		int(width), int(height), sdl.WINDOW_SHOWN)
	if err != nil {
		log.Fatalf("Failed to create window: %s\n", err)
	}
	defer window.Destroy()

	renderer, err = sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED|sdl.RENDERER_PRESENTVSYNC)
	if err != nil {
		log.Fatalf("Failed to create renderer: %s\n", err)
	}
	defer renderer.Destroy()

	if ttf.Init() != 0 {
		log.Fatalf("Failed to init ttf\n")
	}

	font, err = ttf.OpenFont("fonts/arcade_n.ttf", 24)
	if err != nil {
		log.Fatalf("Failed to open font: %s\n", err)
	}

	s := scene.NewScene(200, width, height)

	// three leaders required for default triggers
	s.AddLeader()
	s.AddLeader()
	s.AddLeader()

	// send note off for all the triggers
	Panic(outputStream, s)

	// main loop
	for running {

		for event = sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch t := event.(type) {
			case *sdl.QuitEvent:
				running = false
			case *sdl.KeyUpEvent:
				switch t.Keysym.Sym {
				case sdl.K_q:
					running = false
				case sdl.K_1:
					s.ShowHUD = !s.ShowHUD
				case sdl.K_2:
					s.ShowGrid = !s.ShowGrid
				case sdl.K_3:
					s.ShowActivePads = !s.ShowActivePads
				case sdl.K_UP:
					s.Distance += 0.001
					log.Println("Further apart")
				case sdl.K_DOWN:
					s.Distance -= 0.001
					log.Println("Closer together")
				case sdl.K_RIGHT:
					s.Speed = s.Speed + s.Speed*0.2
					log.Println("Faster")
				case sdl.K_LEFT:
					s.Speed = s.Speed - s.Speed*0.2
					log.Println("Slower")
				case sdl.K_d:
					s.RestoreDefault()
					Panic(outputStream, s)
					log.Println("Set Default Values")
				case sdl.K_l:
					s.AddLeader()
					log.Println("Adding Leader")
				case sdl.K_PERIOD:
					s.RemoveLeader()
					Panic(outputStream, s)
					log.Println("Remove Leader")
				case sdl.K_o:
					s.AllLeaders()
					log.Println("All Leaders")
				case sdl.K_p:
					Panic(outputStream, s)
					log.Println("PANIC")
				}
			}
		}

		s.UpdateBoids()
		events := s.UpdateTriggers()

		if useMidi && len(events) > 0 {
			midiEvents := make([]portmidi.Event, len(events))
			for i, e := range events {
				note := abletonPushNoteMap[e.Note]
				if e.TriggerEventType == scene.TriggerEventType_ON {
					midiEvents[i] = portmidi.Event{
						Timestamp: portmidi.Time(),
						Status:    0x90,
						Data1:     note,
						Data2:     127,
					}
				} else if e.TriggerEventType == scene.TriggerEventType_OFF {
					midiEvents[i] = portmidi.Event{
						Timestamp: portmidi.Time(),
						Status:    0x80,
						Data1:     note,
						Data2:     127,
					}
				}
			}
			outputStream.Write(midiEvents)
		}

		// clear the screen
		renderer.SetDrawColor(0, 0, 0, 255)
		renderer.Clear()

		// render the things
		s.Draw(width, height, renderer, font)
		renderer.Present()
	}

	portmidi.Terminate()
}

func Panic(outputStream *portmidi.Stream, s scene.Scene) {
	// write a note off for all the trigger pads
	if outputStream != nil {

		midiEvents := make([]portmidi.Event, len(s.Triggers))
		for i := range s.Triggers {
			note := abletonPushNoteMap[s.Triggers[i].Note]
			midiEvents[i] = portmidi.Event{
				Timestamp: portmidi.Time(),
				Status:    0x80,
				Data1:     note,
				Data2:     127,
			}

		}
		outputStream.Write(midiEvents)
	}
}

var abletonPushNoteMap = map[int32]int64{
	0:  64,
	1:  65,
	2:  66,
	3:  67,
	4:  96,
	5:  97,
	6:  98,
	7:  99,
	8:  60,
	9:  61,
	10: 62,
	11: 63,
	12: 92,
	13: 93,
	14: 94,
	15: 95,
	16: 56,
	17: 57,
	18: 58,
	19: 59,
	20: 88,
	21: 89,
	22: 90,
	23: 91,
	24: 52,
	25: 53,
	26: 54,
	27: 55,
	28: 84,
	29: 85,
	30: 86,
	31: 87,
	32: 48,
	33: 49,
	34: 50,
	35: 51,
	36: 80,
	37: 81,
	38: 82,
	39: 83,
	40: 44,
	41: 45,
	42: 46,
	43: 47,
	44: 76,
	45: 77,
	46: 78,
	47: 79,
	48: 40,
	49: 41,
	50: 42,
	51: 43,
	52: 72,
	53: 73,
	54: 74,
	55: 75,
	56: 36,
	57: 37,
	58: 38,
	59: 39,
	60: 68,
	61: 69,
	62: 70,
	63: 71}
