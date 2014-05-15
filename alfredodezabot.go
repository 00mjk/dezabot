package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	irc "github.com/thoj/go-ircevent"
)

var rooms []string
var maps map[string]string

func findGif(message string) string {
	for key, value := range maps {
		if strings.Contains(message, key) {
			return value
		}
	}

	return ""
}

func usage() {
	fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
	flag.PrintDefaults()
	os.Exit(1)
}

func main() {
	var (
		userNick    = flag.String("nick", "dezabot", "The IRC nick.")
		userName    = flag.String("user", "dezabot", "The IRC username.")
		serverHost  = flag.String("host", "", "The IRC server host.")
		serverPort  = flag.Int("port", 6667, "The IRC server port.")
		serverPass  = flag.String("pass", "", "The IRC server password.")
		serverTLS   = flag.Bool("usetls", false, "Use TLS for the connection.")
		roomsString = flag.String("rooms", "", "List of rooms to join")
	)

	flag.Parse()

	if *serverHost == "" {
		usage()
	}

	rooms = strings.Split(*roomsString, ",")

	maps = map[string]string{
		"it's magic": "http://media-cache-cd0.pinimg.com/originals/6c/c0/51/6cc0517b7d8aac7515fb75e31865d3b4.jpg",
		"its magic":  "http://media-cache-cd0.pinimg.com/originals/6c/c0/51/6cc0517b7d8aac7515fb75e31865d3b4.jpg",
		"itsmagic":   "http://media-cache-cd0.pinimg.com/originals/6c/c0/51/6cc0517b7d8aac7515fb75e31865d3b4.jpg",
	}

	conn := irc.IRC(*userNick, *userName)
	conn.Password = *serverPass
	conn.UseTLS = *serverTLS

	fmt.Printf("*** Initizlizing dezabot ***\n")

	if err := conn.Connect(fmt.Sprintf("%s:%v", *serverHost, *serverPort)); err != nil {
		log.Fatal(err)
	}

	conn.AddCallback("001", func(e *irc.Event) {
		for _, room := range rooms {
			conn.Join("#" + room)
		}
	})

	conn.AddCallback("PRIVMSG", func(e *irc.Event) {
		if gif := findGif(e.Message()); gif != "" {
			log.Printf("Sending GIF to %s.\n", e.Arguments[0])
			conn.Privmsg(e.Arguments[0], gif)
		}
	})

	conn.Loop()
}
