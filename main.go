package main

import (
	"log"
	"runtime"

	"github.com/veandco/go-sdl2/sdl"
)

var winTitle string = "Der Flock"
var winWidth, winHeight int = 900, 900

func init() {
	runtime.GOMAXPROCS(2)
	runtime.LockOSThread()
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
		winWidth, winHeight, sdl.WINDOW_SHOWN)
	if err != nil {
		log.Fatalf("Failed to create window: %s\n", err)
	}
	defer window.Destroy()

	renderer, err = sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED|sdl.RENDERER_PRESENTVSYNC)
	if err != nil {
		log.Fatalf("Failed to create renderer: %s\n", err)
	}
	defer renderer.Destroy()

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
					// TODO: keyboard actions
				}
			}
		}

		// black background
		renderer.SetDrawColor(0, 0, 0, 255)
		renderer.Clear()

		// render the things
		renderer.SetDrawColor(255, 255, 255, 255)
		renderer.DrawPoint(winHeight/2, winWidth/2)

		renderer.Present()
	}
}
