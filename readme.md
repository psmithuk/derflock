# Der Flock

This is our project for #midihack2 Berlin (May 2015)

## Summary

Artificial birds collaborating to trigger sounds.

It's a MIDI controller based on Boid's Flocking algorithm.
Our Birds fly through a grid of 64 squares triggering sounds in Ableton Live. The movement is also visualised on Ableton Push. You can adjust flocking parameters to change the behaviour of the birds and therefore the sound.

## Details

It's written in Go, using portmidi to sent the MIDI event stream and SDL2 for the visualisation. You should be able to install the dependencies and then `go run`.

It's a one-day hack. Largely uncommented and untested. It will leak memory, slow down and probably destroy your computer. You'll need to change the note mapping table

_We'll add some screenshots, video and sound examples_


