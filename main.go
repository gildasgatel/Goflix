package main

import (
	"goflix/db"
	"goflix/server"
	"log"
)

func main() {

	var db db.Storage = db.New()
	err := db.Setup()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	var server server.Server = server.New(db)
	server.Run()

}
