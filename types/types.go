package types

type Shape struct {
	
	Drawable uint32
	// Values are in pixels and have origin (0,0) at Topleft
	Xpos float32
	Ypos float32
	Width float32
	Height float32

	// List of shape point with OpenGL origin (0,0) at center of the screen
	// and positions are in in OpenGL range (-1,1)
	BasePoints []float32
	Points []float32       // positioned using BasePoints and (Xpos,Ypos)
}

