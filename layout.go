package mcbanner

import "math"

const (
	horizontal = iota
	vertical
)

type Size struct {
	Width  float64
	Height float64
}

type Point struct {
	X float64
	Y float64
}

type TableLayout struct {
	Size     Size
	CellSize Size
}

func (l TableLayout) getLayout(itemCount int) (int, int) {
	if l.getOrientation() == vertical {
		// vertical, try allocate more rows.
		row := math.Floor(l.Size.Height / l.CellSize.Height)
		column := math.Ceil(float64(itemCount) / row)
		return int(row), int(column)
	}
	column := math.Floor(l.Size.Width / l.CellSize.Width)
	row := math.Ceil(float64(itemCount) / column)
	return int(row), int(column)
}

func (l TableLayout) getOrientation() int {
	if l.Size.Width > l.Size.Height {
		return horizontal
	}
	return vertical
}
