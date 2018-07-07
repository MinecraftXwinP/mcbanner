package mcbanner

import (
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"io/ioutil"
	"os"

	"golang.org/x/image/math/fixed"

	"github.com/golang/freetype/truetype"
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

func loadBackground() (draw.Image, error) {
	backgroundFile, err := os.Open("./background.png")
	if err != nil {
		return nil, err
	}
	img, err := png.Decode(backgroundFile)
	if err != nil {
		return nil, err
	}
	v, ok := img.(draw.Image)
	if !ok {
		return nil, fmt.Errorf("%v is not drawable", img)
	}
	return v, nil
}

func loadFontFace() (font.Face, error) {
	fontSrc, err := os.Open("NotoMono-Regular.ttf")
	if err != nil {
		return nil, err
	}
	data, err := ioutil.ReadAll(fontSrc)
	if err != nil {
		return nil, err
	}
	f, err := truetype.Parse(data)
	if err != nil {
		return nil, err
	}
	return truetype.NewFace(f, &truetype.Options{
		Size:    12,
		DPI:     72,
		Hinting: font.HintingFull,
	}), nil
}

func GetDefault(status ServerStatus) (image.Image, error) {
	background, err := loadBackground()
	if err != nil {
		return nil, err
	}
	f, err := loadFontFace()
	if err != nil {
		return nil, err
	}
	// draw adress
	d := &font.Drawer{
		Dst:  background,
		Src:  image.White,
		Face: f,
	}
	rect := background.Bounds()
	width := rect.Dx()
	height := rect.Dy()
	d.Dot = fixed.Point26_6{
		X: fixed.I(width / 15),
		Y: fixed.I(height / 15),
	}
	d.DrawString(fmt.Sprintf("伺服器:%s:%d", status.Host, status.Port))

	return background, nil
}
