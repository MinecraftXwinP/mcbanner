package mcbanner

import (
	"fmt"
	"image"
	"image/draw"
	"math"
	"strconv"

	"golang.org/x/image/font/inconsolata"
	"golang.org/x/image/math/fixed"

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
	MeasureString(s string) (advanced fixed.Int26_6)
}

func (list PlayerList) GetNameWidth(d StringSizeMeasurer) fixed.Int26_6 {
	max := fixed.I(1)
	for _, p := range list.Players {
		w := d.MeasureString(p.Name)
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
	img := image.NewRGBA64(b.Background.Bounds())
	// draw background
	draw.Draw(img, b.Background.Bounds(), b.Background, image.ZP, draw.Src)
	// draw adress
	d := &font.Drawer{
		Dst:  img,
		Src:  image.Black,
		Face: b.FontFace,
	}
	rect := b.Background.Bounds()
	width := rect.Dx()
	height := rect.Dy()
	dx, dy := fixed.I(width/12), fixed.I(height/12)
	d.Dot = fixed.Point26_6{
		X: dx,
		Y: dy,
	}

	d.DrawString(b.ServerStatus.GetAddress())

	countStr := b.ServerStatus.GetPlayerCount()

	// render player count
	d.Dot = fixed.Point26_6{
		X: dx*11 - d.MeasureString(countStr),
		Y: dy,
	}
	d.DrawString(countStr)

	// render player names
	layout := nameListLayout{
		Size: fixed.Point26_6{
			X: dx * 9,
			Y: dy * 9,
		},
		CellSize: fixed.Point26_6{
			X: b.ServerStatus.PlayerList.GetNameWidth(d),
			Y: d.Face.Metrics().Height,
		},
	}

	playerCount := len(b.ServerStatus.PlayerList.Players)
	row, column := layout.getLayout(playerCount)

	p := 0
	for c := 0; c < column; c++ {
		x := dx + layout.CellSize.X.Mul(fixed.I(c))
		for r := 0; r < row; r++ {
			if p >= playerCount-1 {
				break
			}
			d.Dot = fixed.Point26_6{
				X: x,
				Y: dy.Mul(fixed.I(r + 2)),
			}
			d.DrawString(b.ServerStatus.PlayerList.Players[p].Name)
			p++
		}
	}

	return img
}

const (
	horizontal = iota
	vertical
)

type nameListLayout struct {
	Size     fixed.Point26_6
	CellSize fixed.Point26_6
}

func (n nameListLayout) getLayout(itemCount int) (int, int) {
	if n.getOrientation() == vertical {
		// vertical, try allocate more rows.
		row := int(n.Size.Y / n.CellSize.Y)
		column := int(math.Ceil(float64(itemCount) / float64(row)))
		return row, column
	}
	column := int(n.Size.X / n.CellSize.X)
	row := int(math.Ceil(float64(itemCount) / float64(column)))
	return row, column
}

func (n nameListLayout) getOrientation() int {
	if n.Size.X > n.Size.Y {
		return horizontal
	}
	return vertical
}
