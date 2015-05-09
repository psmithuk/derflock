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

	s.Boids[0].BoidKind = scene.BoidKind_LEADER
	s.Boids[1].BoidKind = scene.BoidKind_LEADER
	s.Boids[2].BoidKind = scene.BoidKind_LEADER

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
					// TODO: keyboard actions
				}
			}
		}

		s.UpdateBoids()
		events := s.UpdateTriggers()

		if len(events) > 0 {
			log.Println(events)
		}

		// clear the screen
		renderer.SetDrawColor(0, 0, 0, 255)
		renderer.Clear()

		// render the things
		s.Draw(width, height, renderer)

		renderer.Present()
	}
}
