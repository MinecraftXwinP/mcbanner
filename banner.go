package mcbanner

import (
	"image"
	"image/draw"
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

type PlayerList struct {
	Players []uuid.UUID
	Max     int
	Min     int
}

type ServerStatus struct {
	Host       string
	Port       int
	PlayerList PlayerList
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
	d.Dot = fixed.Point26_6{
		X: fixed.I(width / 15),
		Y: fixed.I(height / 15),
	}
	addr := b.ServerStatus.Host
	if b.ServerStatus.Port != 25565 {
		addr += strconv.Itoa(b.ServerStatus.Port)
	}
	d.DrawString(addr)

	return img
}
