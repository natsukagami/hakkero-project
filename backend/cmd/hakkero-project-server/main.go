package main

import (
	"github.com/natsukagami/hakkero-project/backend"
)

func main() {
	srv := backend.NewServer(backend.DefaultConfig(), backend.StaticOP(), &backend.Rooms{})
	srv.Addr = ":3000"
	go srv.ListenAndServe()
	println("Ready!")
	<-make(chan int)
}
