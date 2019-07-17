package main

import (
	"runtime"
	"github.com/go-gl/glfw/v3.2/glfw"
	grphEngine "gl1/screenwriter"
	gaemEngine "gl1/puppetMaster"
	cfg "gl1/config"
	//"gl1/types"
	"fmt"
)


var (
	ticker = 0
)

const (
	Window_width  = cfg.Window_width
	Window_height = cfg.Window_height
)



func main() {
	// Init
	runtime.LockOSThread()
	window := grphEngine.InitGlfw()
	program := grphEngine.InitOpenGL()
	defer glfw.Terminate()

	fmt.Println(".")
	
	// Create drawable object
	cursor := gaemEngine.NewSquare(10.0, 10.0, 30.0, 30.0)
	//UpdateObjects(cursor)

	// Main loop
	for !window.ShouldClose() {
		// Updates
		cursor.Xpos += 1
		gaemEngine.UpdateShape(cursor)

		// Draw cycle
		grphEngine.ClearScreen(program)
		grphEngine.DrawShape(cursor)
		grphEngine.DrawFrame(window)
	}
}

