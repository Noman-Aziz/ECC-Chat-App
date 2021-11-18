package main

import (
	"flag"
	"math/rand"
	"os"
	"time"

	"github.com/Noman-Aziz/ECC-Chat-App/chat"
)

func main() {

	var username, _ = os.LookupEnv("USER")
	var name = flag.String("name", username, "Name for the client (default $USER)")
	var isHost = flag.Bool("host", false, "Open TCP Port for partner to connect")
	var port = flag.Uint("port", 8080, "Set Port for TCP Socket [For running over network, this port will be forwarded by UPnP]")
	var help = flag.Bool("help", false, "Display Help Page")
	flag.Parse()

	if *help || len(os.Args) < 2 {
		flag.Usage()
		return
	}

	rand.Seed(time.Now().UTC().UnixNano())

	var appConfig = chat.Config{
		Name:   *name,
		IsHost: *isHost,
		Port:   uint16(*port),
	}
	var app = chat.CreateChatApp(&appConfig)

	app.Run()
}
