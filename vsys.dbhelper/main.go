package main

import (
	"log"

	db "vsys.dbhelper/db"
	srv "vsys.dbhelper/server"
)

func main() {

	// Initialize database connection
	database := db.GetDB()
	if database == nil {
		log.Fatal("Failed to establish database connection.")
	}

	// Start the web server
	srv.Web()
}
