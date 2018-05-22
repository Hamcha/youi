package components

import (
	"fmt"
	"image"
)

type Position struct {
	X, Y float32
}

func (p Position) Scale(s Size) Position {
	return Position{p.X * s.Width, p.Y * s.Height}
}

func (p Position) String() string {
	return fmt.Sprintf("(%.2f; %.2f)", p.X, p.Y)
}

type Size struct {
	Width, Height float32
}

func (s Size) String() string {
	return fmt.Sprintf("(%.2f; %.2f)", s.Width, s.Height)
}

func (s Size) Scale(z Size) Size {
	return Size{s.Width * z.Width, s.Height * z.Height}
}

func (s Size) Inverse() Size {
	return Size{1.0 / s.Width, 1.0 / s.Height}
}

type Bounds struct {
	Position
	Size
}

func BoundsFromRect(rect image.Rectangle) Bounds {
	size := rect.Size()
	return Bounds{
		Position{
			X: float32(rect.Min.X),
			Y: float32(rect.Min.Y),
		},
		Size{
			Width:  float32(size.X),
			Height: float32(size.Y),
		},
	}
}

func (b Bounds) Scale(s Size) Bounds {
	return Bounds{
		Position: b.Position.Scale(s),
		Size:     b.Size.Scale(s),
	}
}

func (b Bounds) String() string {
	return fmt.Sprintf("Position %s Size %s", b.Position, b.Size)
}
