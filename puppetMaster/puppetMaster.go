package puppetMaster

import (
	grphEngine "gl1/screenwriter"
	"gl1/types"
)

func NewSquare(Width_, Height_, Xpos_, Ypos_ float32) *types.Shape{
	// Create new shape struct
	obj := &types.Shape{
		Width: Width_,
		Height: Height_,
		Xpos: Xpos_,
		Ypos: Ypos_,
	}

	// Create a base square with a certain width and height
	obj.BasePoints = grphEngine.RenderSquarePoints(obj.Width, obj.Height)

	// Translate base square to final square and create vao
	UpdateShape(obj)
	
	return obj
}

func UpdateShape(obj *types.Shape){
	// Calculate obj.Points from obj.BasePoints and its position
	grphEngine.PositionShape(obj)

	// Create vao from obj.Points
	obj.Drawable = grphEngine.MakeVao(obj.Points)
}