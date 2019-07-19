package main

import (
	"runtime"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/go-gl/gl/v2.1/gl"
	grphEngine "gl1/screenwriter"
	gaemEngine "gl1/puppetMaster"
	cfg "gl1/config"
	"gl1/types"
	"fmt"
	"time"
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

	grphEngine.InitVBO()

	fmt.Println(".")

	// Create drawable objects
	//squares := make([]*types.Shape, 7000)
	//squares := make(map[int]*types.Shape)
	xoffset := -37
	yoffset := 50

	var current, previous, firstSquare *types.Shape
	numberof_active_squares := 0

	for i := 0; i < 20; i++ {
		switch i % 10 {
			case 0:
				yoffset = 50
				xoffset += 37
			default:
				yoffset += 22
				 
		}
		// new square
		current = gaemEngine.NewSquare(10.0, 10.0, float32(-xoffset), float32(yoffset), i)
		
		// if there is no previous, this is the first square in the linked list
		if previous == nil {
			firstSquare = current
		} else {
			// there is a previous square, let eachother know the current and previous
			previous.NextShape = current
			current.PreviousShape = previous
		}
		// switch context in preparation for next iteration
		previous = current
	}

	/* Main loop -------------------------------------------------------- */
	for !window.ShouldClose() {
		start := time.Now()

		/* Clear screen -------------------------------------------------------- */
		grphEngine.ClearScreen(program)

		/* Updates -------------------------------------------------------- */

		// ## Loop over all squares
		current = firstSquare
		numberof_active_squares = 0
		for {
			// End of linked list reached
			if current == nil{
				break
			}
			// # Update position
			current.SetPos(current.Xpos + 1, current.Ypos)

			// # Draw square
			if current.OnScreen {
				switch numberof_active_squares % 2{
					case 0:

					default:

				}
				gaemEngine.UpdateShape(current)
				grphEngine.DrawShape(current)
				numberof_active_squares += 1
			}
			// # Kill square
			if current.Kill {
				gl.DeleteVertexArrays(1, &(current.Drawable))

				if current.PreviousShape == nil    {	
					// firstSquare is killed, set nextshape to be firstshape	
					firstSquare = current.NextShape    
					if firstSquare != nil{	
						firstSquare.PreviousShape = nil
					}

				} else if current.NextShape == nil { 			
					// lastSquare is killed, set previousshape to be lastsquare
					(*current.PreviousShape).NextShape = nil

				} else {										
					(*current.PreviousShape).NextShape = current.NextShape
					(*current.NextShape).PreviousShape = current.PreviousShape
				}
			}
			// # Switch context
			current = current.NextShape
		}

		/* Draw new frame -------------------------------------------------------- */
		grphEngine.DrawFrame(window)

		
		t := time.Now()
		elapsed := t.Sub(start)
		fmt.Println(1E9/elapsed.Nanoseconds())
		

		//fmt.Println(numberof_active_squares)
	}
}

