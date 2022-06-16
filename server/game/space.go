package game

type Point struct {
	x, y int
}

func NewPoint(x, y int) *Point {
	return &Point{
		x: x,
		y: y,
	}
}
