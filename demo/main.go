package main

import (
	"flag"
	"fmt"
	"image"
	"image/png"
	"io/ioutil"
	"log"
	"net"
	"os"
	"strconv"

	"github.com/golang/freetype/truetype"
	"github.com/google/uuid"

	"github.com/andrewtian/minepong"

	"github.com/MinecraftXwinP/mcbanner"
)

var host string
var port int
var backgroundPath string

func init() {
	flag.StringVar(&host, "server", "", "minecraft server address")
	flag.StringVar(&backgroundPath, "bkg", "", "custom background image. format should be PNG.")
	flag.IntVar(&port, "port", 25565, "server port.")
	flag.Parse()
}

func main() {
	if len(host) <= 0 {
		fmt.Printf("host: %s. port: %d. bkg: %s", host, port, backgroundPath)
		flag.Usage()
		os.Exit(-1)
	}
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
	status := mcbanner.ServerStatus{
		Host: host,
		Port: port,
		PlayerList: mcbanner.PlayerList{
			Max:     pong.Players.Max,
			Players: players,
		},
	}

	var img image.Image
	if len(backgroundPath) <= 0 {
		img = mcbanner.GetDefault(status)
	} else {
		fontFile, err := os.Open("NotoSansDisplay-Regular.ttf")
		if err != nil {
			log.Fatal(err)
		}
		defer fontFile.Close()
		data, err := ioutil.ReadAll(fontFile)
		if err != nil {
			log.Fatal(err)
		}
		f, err := truetype.Parse(data)
		bkFile, err := os.Open(backgroundPath)
		if err != nil {
			log.Fatal(err)
		}
		defer bkFile.Close()
		b, err := png.Decode(bkFile)
		r := mcbanner.Banner{
			Background: b,
			FontFace: truetype.NewFace(f, &truetype.Options{
				Size: 16,
			}),
			ServerStatus: status,
		}
		img = r.Render()
	}
	out, err := os.Create("./out.png")
	if err != nil {
		log.Fatal(err)
	}
	err = png.Encode(out, img)
	if err != nil {
		log.Fatal(err)
	}
}
