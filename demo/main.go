package main

import (
	"fmt"
	"image/png"
	"log"
	"os"

	"github.com/google/uuid"

	"github.com/MinecraftXwinP/mcbanner"
)

func main() {
	fakePlayers := make([]mcbanner.Player, 20)
	for i := range fakePlayers {
		fakePlayers[i] = mcbanner.Player{
			Name: fmt.Sprintf("Player%d", i),
			UUID: uuid.New(),
		}
	}
	img := mcbanner.GetDefault(mcbanner.ServerStatus{
		Host: "example.com",
		Port: 25565,
		PlayerList: mcbanner.PlayerList{
			Max:     20,
			Min:     0,
			Players: fakePlayers,
		},
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
