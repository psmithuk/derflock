package main

import (
	"log"
	"math/rand"
	"runtime"
	"time"

	"github.com/psmithuk/derflock/scene"
	"github.com/rakyll/portmidi"
	"github.com/veandco/go-sdl2/sdl"
)

var winTitle string = "Der Flock"
var width, height int32 = 720, 720

func init() {
	runtime.GOMAXPROCS(2)
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
			log.Fatalf("Failed open new midi output: %s\n", err)
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

	s := scene.NewScene(200, width, height)

	s.AddLeader()
	s.AddLeader()
	s.AddLeader()

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
					log.Println("Set Default Values")
				case sdl.K_l:
					s.AddLeader()
					log.Println("Adding Leader")
				case sdl.K_o:
					s.AllLeaders()
					log.Println("All Leaders")
					// TODO: keyboard actions
				}
			}
		}

		s.UpdateBoids()
		events := s.UpdateTriggers()

		if useMidi && len(events) > 0 {
			for _, e := range events {
				if e.TriggerEventType == scene.TriggerEventType_ON {
					outputStream.WriteShort(0x90, int64(e.Note), 100)
				} else if e.TriggerEventType == scene.TriggerEventType_OFF {
					outputStream.WriteShort(0x80, int64(e.Note), 100)
				}
			}
		}

		// clear the screen
		renderer.SetDrawColor(0, 0, 0, 255)
		renderer.Clear()

		// render the things
		s.Draw(width, height, renderer)

		renderer.Present()
	}
}
