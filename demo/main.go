package main

import (
	"fmt"
	"image/png"
	"log"
	"net"
	"os"
	"strconv"

	"github.com/google/uuid"

	"github.com/andrewtian/minepong"

	"github.com/MinecraftXwinP/mcbanner"
)

var host string
var port int

func init() {
	argc := len(os.Args)
	if argc < 2 {
		fmt.Println("Usage: demo [server address] [port]")
		os.Exit(-1)
	}
	host = os.Args[1]
	var err error
	if argc > 2 {
		port, err = strconv.Atoi(os.Args[2])
		if err != nil {
			log.Fatal("invalid port", err)
		}
	} else {
		port = 25565
	}
}

func main() {
	addr := host + ":" + strconv.Itoa(port)
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		log.Fatal("cannot connect to the server", err)
	}
	pong, err := minepong.Ping(conn, addr)
	if err != nil {
		log.Fatal("cannot ping the server", err)
	}
	players := make([]mcbanner.Player, len(pong.Players.Sample))
	for i, p := range pong.Players.Sample {
		id, err := uuid.Parse(p["id"])
		if err != nil {
			log.Fatal("fail to parse player id", err)
		}
		players[i] = mcbanner.Player{
			Name: p["name"],
			UUID: id,
		}
	}

	img := mcbanner.GetDefault(mcbanner.ServerStatus{
		Host: host,
		Port: port,
		PlayerList: mcbanner.PlayerList{
			Max:     pong.Players.Max,
			Players: players,
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
