package mcbanner

import (
	"fmt"
	"image"
	"image/draw"
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
	c := newDrawerForImage(b.Background)
	c.DrawImage(b.Background, 0, 0)
	c.SetRGB255(0, 0, 0)
	c.SetFontFace(b.FontFace)
	c.DrawOutlinedTextAnchored(b.ServerStatus.GetAddress(), float64(0), float64(0), 0, 1)
	c.DrawOutlinedTextAnchored(b.ServerStatus.GetPlayerCount(), float64(width), float64(0), 1, 1)

	_, textHeight := c.MeasureString("A")

	layout := TableLayout{
		Size: Size{
			float64(width),
			float64(height * 2 / 3),
		},
		CellSize: Size{
			b.ServerStatus.PlayerList.GetNameWidth(c) + 30,
			textHeight * 1.2,
		},
	}

	playerCount := len(b.ServerStatus.PlayerList.Players)
	row, column := layout.getLayout(playerCount)

	c.Push()
	p := 0
	for col := 0; col < column; col++ {
		x := layout.CellSize.Width * float64(col)
		for r := 0; r < row; r++ {
			if p >= playerCount {
				break
			}
			c.Push()
			y := float64(r)*layout.CellSize.Height + (float64(height) / 3)
			c.SetRGBA255(11, 11, 11, 80)
			c.DrawRectangle(x, y, layout.CellSize.Width, layout.CellSize.Height)
			c.Fill()
			c.Pop()

			c.DrawOutlinedTextAnchored(b.ServerStatus.PlayerList.Players[p].Name, x, y, 0, 1)
			p++
		}
	}

	return c.Image()
}

type drawer struct {
	*gg.Context
}

func (d *drawer) DrawOutlinedText(message string, x, y float64) {
	d.Push()
	d.SetRGB(0, 0, 0)
	d.DrawString(message, x-1, y-1)
	d.DrawString(message, x+1, y+1)
	d.Pop()
	d.Push()
	d.SetRGB(1, 1, 1)
	d.DrawString(message, x, y)
	d.Pop()
}

func (d *drawer) DrawOutlinedTextAnchored(message string, x, y, ax, ay float64) {
	d.Push()
	d.SetRGB(0, 0, 0)
	d.DrawStringAnchored(message, x-1, y-1, ax, ay)
	d.DrawStringAnchored(message, x+1, y+1, ax, ay)
	d.Pop()
	d.Push()
	d.SetRGB(1, 1, 1)
	d.DrawStringAnchored(message, x, y, ax, ay)
	d.Pop()
}

func newDrawer(width, height int) *drawer {
	return &drawer{
		gg.NewContext(width, height),
	}
}

func newDrawerForImage(img image.Image) *drawer {
	return &drawer{
		gg.NewContextForImage(img),
	}
}
