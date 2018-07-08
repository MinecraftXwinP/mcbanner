package mcbanner

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"golang.org/x/image/math/fixed"
)

func TestGetOrientation(t *testing.T) {
	layout := nameListLayout{
		Size: fixed.Point26_6{
			X: fixed.I(200),
			Y: fixed.I(20),
		},
	}
	assert.Equal(t, horizontal, layout.getOrientation())
}

func TestHorizontalLayout(t *testing.T) {
	layout := nameListLayout{
		Size:     fixed.P(400, 200),
		CellSize: fixed.P(20, 5),
	}

	row, column := layout.getLayout(100)
	assert.Equal(t, 20, column)
	assert.Equal(t, 5, row)
}

func TestHorizontalLayoutWithFraction(t *testing.T) {
	layout := nameListLayout{
		Size:     fixed.P(510, 200),
		CellSize: fixed.P(20, 5),
	}
	row, column := layout.getLayout(99)
	assert.Equal(t, 25, column)
	assert.Equal(t, 4, row)
}

func TestVerticalLayout(t *testing.T) {
	layout := nameListLayout{
		Size:     fixed.P(200, 400),
		CellSize: fixed.P(20, 5),
	}

	row, column := layout.getLayout(160)
	assert.Equal(t, 80, row)
	assert.Equal(t, 2, column)
}

func TestVerticalLayoutWithFraction(t *testing.T) {
	layout := nameListLayout{
		Size:     fixed.P(199, 400),
		CellSize: fixed.P(20, 7),
	}
	row, column := layout.getLayout(199)
	assert.Equal(t, 57, row)
	assert.Equal(t, 4, column)

}

func ExampleGetAddress() {
	status := ServerStatus{
		Host: "example.com",
		Port: 25565,
	}
	fmt.Println(status.GetAddress())

	status = ServerStatus{
		Host: "example.com",
		Port: 25566,
	}

	fmt.Println(status.GetAddress())
	// Output:
	// example.com
	// example.com:25566
}

// font.Drawer Mock
type fontDrawerMock struct {
}

func (fd *fontDrawerMock) MeasureString(s string) fixed.Int26_6 {
	return fixed.I(len(s))
}

func TestGetNameWidth(t *testing.T) {
	fakePlayers := make([]Player, 20)
	for i := 0; i < 19; i++ {
		fakePlayers[i] = Player{
			Name: fmt.Sprintf("Player%d", i),
			UUID: uuid.New(),
		}
	}
	longName := "longest player name    "
	fakePlayers[19] = Player{
		Name: longName,
		UUID: uuid.New(),
	}
	list := PlayerList{
		Max:     20,
		Players: fakePlayers,
	}
	assert.Equal(t, fixed.I(len(longName)), list.GetNameWidth(&fontDrawerMock{}))
}

func TestGetNameWidthWillNotReturnZero(t *testing.T) {
	namelessPlayers := make([]Player, 20)
	for i := 0; i < 20; i++ {
		namelessPlayers[i] = Player{
			Name: "",
			UUID: uuid.New(),
		}
	}
	list := PlayerList{
		Max:     20,
		Players: namelessPlayers,
	}

	assert.Equal(t, fixed.I(1), list.GetNameWidth(&fontDrawerMock{}))
}
