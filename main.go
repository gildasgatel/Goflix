package main

import (
	"goflix/db"
	"goflix/server"
)

func main() {

	var db db.Storage = db.New()
	db.Setup()
	defer db.Close()

	var server server.Server = server.New(db)
	server.Run()

}
