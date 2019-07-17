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
)

func RenderSquarePoints(Width, Height float32) []float32 {
	// pixel definition
	sqWidth := Width
	sqHeight := Height

	// gl sizes
	sqW := sqWidth * xTransl
	sqH := sqHeight * yTransl

	x_origin :=  float32(0)
	y_origin := float32(0)

	sqPoints := make([]float32,18)
	// topleft
	sqPoints[0] = x_origin
	sqPoints[1] = y_origin
	sqPoints[2] = 0

	// topright
	sqPoints[3] = x_origin + sqW
	sqPoints[4] = y_origin
	sqPoints[5] = 0

	// bottomleft
	sqPoints[6] = x_origin
	sqPoints[7] = y_origin - sqH
	sqPoints[8] = 0

	// topright
	sqPoints[9]  = x_origin + sqW
	sqPoints[10] = y_origin
	sqPoints[11] = 0

	// bottomleft
	sqPoints[12] = x_origin
	sqPoints[13] = y_origin - sqH
	sqPoints[14] = 0

	// bottomright
	sqPoints[15] = x_origin + sqW
	sqPoints[16] = y_origin - sqH
	sqPoints[17] = 0
	
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
	// write to buffer
	gl.BindVertexArray(shape.Drawable)
	gl.DrawArrays(gl.TRIANGLES, 0, int32(len(shape.Points)/3))
}

func DrawFrame(window *glfw.Window) {
	glfw.PollEvents()
	window.SwapBuffers()
}

// makeVao initializes and returns a vertex array from the points provided.
func MakeVao(points []float32) uint32 {
	var vbo uint32
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(points), gl.Ptr(points), gl.STATIC_DRAW)

	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)
	gl.EnableVertexAttribArray(0)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 0, nil)

	return vao
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

	fragmentShader, err := compileShader(fragmentShaderSource, gl.FRAGMENT_SHADER)
	if err != nil {
		panic(err)
	}

	prog := gl.CreateProgram()
	gl.AttachShader(prog, vertexShader)
	gl.AttachShader(prog, fragmentShader)
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
