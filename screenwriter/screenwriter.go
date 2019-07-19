package screenwriter

import (
	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"gl1/types"
	cfg "gl1/config"
	"log"
	"strings"
	"fmt"
)

var (
	Vbo1 uint32
	Vbo2 uint32
	FragmentShader1 uint32
	FragmentShader2 uint32
)
const (
	Window_width  = cfg.Window_width
	Window_height = cfg.Window_height

	xTransl = float32(2.0/cfg.Window_width)
	yTransl = float32(2.0/cfg.Window_height)

	vertexShaderSource = `
		#version 120
		void main() {
		gl_Position = gl_ProjectionMatrix * gl_ModelViewMatrix * gl_Vertex;
		}
	` + "\x00"

	fragmentShaderSource = `
		#version 120
		void main() {
		gl_FragColor = vec4(1, 0.5, 1, 1.0);
		}
	` + "\x00"

	fragmentShaderSourceRed = `
		#version 120
		void main() {
		gl_FragColor = vec4(1, 0, 0, 1.0);
		}
	` + "\x00"

	fragmentShaderSourceGreen = `
		#version 120
		void main() {
		gl_FragColor = vec4(1, 0, 0, 1.0);
		}
	` + "\x00"	
)

func InitVBO(){
	gl.GenBuffers(1, &Vbo1) //red
	gl.GenBuffers(1, &Vbo2) //green
}

func RenderSquarePoints(Width, Height float32) []float32 {
	// pixel definition
	sqWidth := Width
	sqHeight := Height

	// gl sizes
	sqW := sqWidth * xTransl
	sqH := sqHeight * yTransl

	// init
	x_origin := float32(0)
	y_origin := float32(0)
	origin := []float32{x_origin, y_origin, 0}; 

	var sqPoints []float32
	
	// define corners
	topleft     := make([]float32,3); copy(topleft,     origin)
	topright    := make([]float32,3); copy(topright,    origin)
	bottomleft  := make([]float32,3); copy(bottomleft,  origin)
	bottomright := make([]float32,3); copy(bottomright, origin)
	
	topright[0]    += sqW
	bottomleft[1]  -= sqH
	bottomright[0] += sqW
	bottomright[1] -= sqH

	// triangle 1
	sqPoints = append(sqPoints, topleft...)
	sqPoints = append(sqPoints, topright...)
	sqPoints = append(sqPoints, bottomleft...)

	// triangle 2
	sqPoints = append(sqPoints, topright...)
	sqPoints = append(sqPoints, bottomleft...)
	sqPoints = append(sqPoints, bottomright...)

	return sqPoints
}

func PositionShape(shape *types.Shape) {
	// Translate pixel position to OpenGL format
	yPos := shape.Ypos * yTransl
	xPos := shape.Xpos * xTransl

	// Create new Slice with updated points
	shape.Points = make([]float32, len(shape.BasePoints))
	for i := 0; i < len(shape.BasePoints); i++ {
		switch i % 3 {
		// x
		case 0:
			shape.Points[i] = -(1.01 - (xPos)) + shape.BasePoints[i]
		// y
		case 1:
			shape.Points[i] =  (1.00 - (yPos)) + shape.BasePoints[i]
		// z
		default:
			continue
		}
	}
}

func ClearScreen(program uint32) {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	gl.UseProgram(program)
}

func DrawShape(shape *types.Shape) {
	if shape.OnScreen {
		// write to buffer
		gl.BindVertexArray(shape.Drawable)
		gl.DrawArrays(gl.TRIANGLES, 0, int32(len(shape.Points)/3))
	}
}

func DrawFrame(window *glfw.Window) {
	glfw.PollEvents()
	window.SwapBuffers()
}

// makeVao initializes and returns a vertex array from the points provided.
func MakeVao(points []float32, vao *uint32, color int) uint32 {

	if color == 1 {
		gl.BindBuffer(gl.ARRAY_BUFFER, Vbo1)
	} else {
		gl.BindBuffer(gl.ARRAY_BUFFER, Vbo2)
	}
		
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(points), gl.Ptr(points), gl.STATIC_DRAW)

	if *vao == uint32(0){
		gl.GenVertexArrays(1, vao)
	}

	//gl.BindVertexArray(*vao)
	gl.EnableVertexAttribArray(0)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 0, nil)

	return *vao
}

// initGlfw initializes glfw and returns a Window to use.
func InitGlfw() *glfw.Window {
	if err := glfw.Init(); err != nil {
		panic(err)
	}
	glfw.WindowHint(glfw.Resizable, glfw.False)

	window, err := glfw.CreateWindow(Window_width, Window_height, "Graphatroner", nil, nil)
	if err != nil {
		panic(err)
	}
	window.MakeContextCurrent()

	return window
}

// initOpenGL initializes OpenGL and returns an intiialized program.
func InitOpenGL() uint32 {
	if err := gl.Init(); err != nil {
		panic(err)
	}
	version := gl.GoStr(gl.GetString(gl.VERSION))
	log.Println("OpenGL version", version)

	vertexShader, err := compileShader(vertexShaderSource, gl.VERTEX_SHADER)
	if err != nil {
		panic(err)
	}

	FragmentShader, err := compileShader(fragmentShaderSource, gl.FRAGMENT_SHADER)
	if err != nil {
		panic(err)
	}
	
	prog := gl.CreateProgram()
	gl.AttachShader(prog, vertexShader)
	gl.AttachShader(prog, FragmentShader)
	gl.LinkProgram(prog)
	return prog
}

func compileShader(source string, shaderType uint32) (uint32, error) {
	shader := gl.CreateShader(shaderType)

	csources, free := gl.Strs(source)
	gl.ShaderSource(shader, 1, csources, nil)
	free()
	gl.CompileShader(shader)

	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(log))

		return 0, fmt.Errorf("failed to compile %v: %v", source, log)
	}

	return shader, nil
}
