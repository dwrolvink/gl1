package types

import (
	cfg "gl1/config"
)

type Shape struct {
	
	Drawable uint32
	// Values are in pixels and have origin (0,0) at Topleft
	Xpos float32
	Ypos float32
	Width float32
	Height float32
	Changed bool
	OnScreen bool
	Name int
	Kill bool
	NextShape *Shape
	PreviousShape *Shape

	// List of shape point with OpenGL origin (0,0) at center of the screen
	// and positions are in in OpenGL range (-1,1)
	BasePoints []float32
	Points []float32       // positioned using BasePoints and (Xpos,Ypos)

}
func (c *Shape) SetPos(Xpos float32, Ypos float32) {
	c.Xpos = Xpos; c.Ypos = Ypos
	c.Changed = true
	
	// Set OnScreen toggle [entering screen]
	if c.Xpos <= 0 {
		c.OnScreen = false
	} else {
		c.OnScreen = true
	}

	// Set Onscreen toggle [leaving screen]
	if c.Xpos >= cfg.Window_width {
		c.Kill = true
	}
}
