package main

import (
	"time"

	"github.com/natsukagami/hakkero-project/backend"
)

func main() {
	srv := backend.NewServer(backend.Config{PlayerLimit: 2, Timeout: 10 * time.Second}, backend.StaticOP(), &backend.Rooms{})
	srv.Addr = ":2020"
	println("Ready!")
	err := srv.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
