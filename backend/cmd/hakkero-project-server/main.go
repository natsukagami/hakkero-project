package main

import (
	"flag"
	"time"

	"github.com/natsukagami/hakkero-project/backend"
)

var (
	playerLimit = flag.Int("player", 4, "Set the player limit on each room.")
	timeout     = flag.Int("timeout", 60, "Set the timeout of each turn")
)

func main() {
	flag.Parse()
	srv := backend.NewServer(backend.Config{PlayerLimit: *playerLimit, Timeout: time.Duration(*timeout) * time.Second}, backend.StaticOP(), &backend.Rooms{})
	srv.Addr = ":80"
	println("Ready!")
	err := srv.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
