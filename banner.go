package mcbanner

import (
	"fmt"
	"image"
	"image/draw"
	"math"
	"strconv"

	"github.com/fogleman/gg"

	"golang.org/x/image/font/inconsolata"

	"github.com/google/uuid"
	"golang.org/x/image/font"
)

const (
	DefaultWidth  = 336
	DefaultHeight = 280
)

type Player struct {
	Name string
	UUID uuid.UUID
}

type PlayerList struct {
	Players []Player
	Max     int
	Min     int
}
type StringSizeMeasurer interface {
	MeasureString(s string) (w, h float64)
}

func (list PlayerList) GetNameWidth(d StringSizeMeasurer) float64 {
	max := float64(1)
	for _, p := range list.Players {
		w, _ := d.MeasureString(p.Name)
		if w > max {
			max = w
		}
	}
	return max
}

type ServerStatus struct {
	Host       string
	Port       int
	PlayerList PlayerList
}

// GetAddress returns server address. If the port is default minecraft server port, 25565, the tcp port will not be shown.
func (status *ServerStatus) GetAddress() string {
	addr := status.Host
	if status.Port != 25565 {
		addr += ":" + strconv.Itoa(status.Port)
	}
	return addr
}

func (status *ServerStatus) GetPlayerCount() string {
	return fmt.Sprintf("%d / %d", len(status.PlayerList.Players), status.PlayerList.Max)
}

// GetDefault returns default banner. White background and with default size 336 x 280
func GetDefault(status ServerStatus) image.Image {
	white := image.NewRGBA64(image.Rect(0, 0, DefaultWidth, DefaultHeight))
	draw.Draw(white, white.Bounds(), image.White, image.ZP, draw.Src)
	b := Banner{
		Background:   white,
		FontFace:     inconsolata.Bold8x16,
		ServerStatus: status,
	}
	return b.Render()
}

type Banner struct {
	ServerStatus ServerStatus
	Background   image.Image
	FontFace     font.Face
}

func (b *Banner) Render() image.Image {
	bounds := b.Background.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()
	c := gg.NewContextForImage(b.Background)
	c.DrawImage(b.Background, 0, 0)
	c.SetRGB255(0, 0, 0)
	c.SetFontFace(b.FontFace)
	c.DrawStringAnchored(b.ServerStatus.GetAddress(), float64(0), float64(0), 0, 1)
	c.DrawStringAnchored(b.ServerStatus.GetPlayerCount(), float64(width), float64(0), 1, 1)

	_, textHeight := c.MeasureString("A")

	layout := nameListLayout{
		Size: Size{
			float64(width),
			float64(height * 2 / 3),
		},
		CellSize: Size{
			b.ServerStatus.PlayerList.GetNameWidth(c) + 30,
			textHeight,
		},
	}

	playerCount := len(b.ServerStatus.PlayerList.Players)
	row, column := layout.getLayout(playerCount)

	p := 0
	for col := 0; col < column; col++ {
		x := layout.CellSize.Width * float64(col)
		for r := 0; r < row; r++ {
			if p >= playerCount {
				break
			}
			// d.DrawString(b.ServerStatus.PlayerList.Players[p].Name)
			c.DrawStringAnchored(b.ServerStatus.PlayerList.Players[p].Name, x, float64(r)*layout.CellSize.Height+(float64(height)/3), 0, 1)
			p++
		}
	}

	return c.Image()
}

const (
	horizontal = iota
	vertical
)

type Size struct {
	Width  float64
	Height float64
}

type nameListLayout struct {
	Size     Size
	CellSize Size
}

func (n nameListLayout) getLayout(itemCount int) (int, int) {
	if n.getOrientation() == vertical {
		// vertical, try allocate more rows.
		row := math.Floor(n.Size.Height / n.CellSize.Height)
		column := math.Ceil(float64(itemCount) / row)
		return int(row), int(column)
	}
	column := math.Floor(n.Size.Width / n.CellSize.Width)
	row := math.Ceil(float64(itemCount) / column)
	return int(row), int(column)
}

func (n nameListLayout) getOrientation() int {
	if n.Size.Width > n.Size.Height {
		return horizontal
	}
	return vertical
}
