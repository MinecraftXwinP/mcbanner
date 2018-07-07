package mcbanner

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"io/ioutil"
	"os"

	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"github.com/google/uuid"
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

func loadFontConctext() (*freetype.Context, error) {
	fontSrc, err := os.Open("NotoMono-Regular.ttf")
	if err != nil {
		return nil, err
	}
	data, err := ioutil.ReadAll(fontSrc)
	if err != nil {
		return nil, err
	}
	font, err := truetype.Parse(data)
	if err != nil {
		return nil, err
	}
	c := freetype.NewContext()
	c.SetFont(font)
	return c, nil
}

func GetDefault(status ServerStatus) (image.Image, error) {
	background, err := loadBackground()
	if err != nil {
		return nil, err
	}
	fontC, err := loadFontConctext()
	if err != nil {
		return nil, err
	}
	// draw address
	fontC.SetDst(background)
	fontC.SetSrc(image.NewUniform(color.RGBA{200, 100, 0, 255}))
	pt := freetype.Pt(20, 20+int(fontC.PointToFixed(14)>>6))
	fontC.DrawString(status.Host, pt)

	return background, nil
}
