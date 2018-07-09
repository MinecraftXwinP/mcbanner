package mcbanner

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func S(width, height float64) Size {
	return Size{
		width,
		height,
	}
}

func TestGetOrientation(t *testing.T) {
	layout := TableLayout{
		Size: S(200, 20),
	}
	assert.Equal(t, horizontal, layout.getOrientation())
}

func TestHorizontalLayout(t *testing.T) {
	layout := TableLayout{
		Size:     S(400, 200),
		CellSize: S(20, 5),
	}

	row, column := layout.getLayout(100)
	assert.Equal(t, 20, column)
	assert.Equal(t, 5, row)
}

func TestHorizontalLayoutWithFraction(t *testing.T) {
	layout := TableLayout{
		Size:     S(510, 200),
		CellSize: S(20, 5),
	}
	row, column := layout.getLayout(99)
	assert.Equal(t, 25, column)
	assert.Equal(t, 4, row)
}

func TestVerticalLayout(t *testing.T) {
	layout := TableLayout{
		Size:     S(200, 400),
		CellSize: S(20, 5),
	}

	row, column := layout.getLayout(160)
	assert.Equal(t, 80, row)
	assert.Equal(t, 2, column)
}

func TestVerticalLayoutWithFraction(t *testing.T) {
	layout := TableLayout{
		Size:     S(199, 400),
		CellSize: S(20, 7),
	}
	row, column := layout.getLayout(199)
	assert.Equal(t, 57, row)
	assert.Equal(t, 4, column)

}
