package mcbanner

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

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

func (fd *fontDrawerMock) MeasureString(s string) (width, height float64) {
	return float64(len(s)), 0
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
	assert.Equal(t, float64(len(longName)), list.GetNameWidth(&fontDrawerMock{}))
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

	assert.Equal(t, float64(1), list.GetNameWidth(&fontDrawerMock{}))
}
