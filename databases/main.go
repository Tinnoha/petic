package main

import database "databases/database"

func main() {
	db := database.NewDatabase()
	database.StartKafka(db)
}
