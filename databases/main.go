package main

import (
	"databases/htttp"
	database "databases/htttp/database"
)

func main() {
	db := database.NewDatabase()
	db.Start()
	hh := htttp.NewHttpHandler(db)
	hs := htttp.NewHttpServer(*hh)

	hs.Run()
}
