package graphics

import rl "github.com/gen2brain/raylib-go/raylib"

type CircularLoader struct {
	Pos         rl.Vector2
	LoaderColor rl.Color
	BorderColor rl.Color
	Radius      float32
	Thickness   float32
	ArcSize     float32
}

func NewCircularLoader(x, y, radius float32) *CircularLoader {
	return &CircularLoader{
		Pos:         rl.NewVector2(x, y),
		LoaderColor: rl.SkyBlue,
		BorderColor: rl.DarkBlue,
		Radius:      radius,
		Thickness:   10,
		ArcSize:     0,
	}
}

func (c *CircularLoader) Update(angle float32) {
	c.ArcSize = angle
}

func (c *CircularLoader) Render() {

	innerRadius := c.Radius - c.Thickness
	outerRadius := c.Radius

	// Background full ring
	rl.DrawRing(
		c.Pos,
		innerRadius,
		outerRadius,
		0,
		360,
		100,
		c.BorderColor,
	)

	// Rotating arc
	startAngle := float32(0)
	endAngle := c.ArcSize

	rl.DrawRing(
		c.Pos,
		innerRadius,
		outerRadius,
		startAngle,
		endAngle,
		100,
		c.LoaderColor,
	)
}
