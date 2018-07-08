package main

import (
	"image/png"
	"log"
	"os"

	"github.com/MinecraftXwinP/mcbanner"
)

func main() {
	img := mcbanner.GetDefault(mcbanner.ServerStatus{
		Host: "example.com",
		Port: 25565,
	})
	out, err := os.Create("./out.png")
	if err != nil {
		log.Fatal(err)
	}
	err = png.Encode(out, img)
	if err != nil {
		log.Fatal(err)
	}
}
