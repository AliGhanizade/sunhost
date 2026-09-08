package main

import (
	"log"
	"os"
	"sunhost/config"
	r "sunhost/router"
)

func main() {
	config.InitDB()
	config.MigrateUser()
	config.MigrateSystemStat()
	config.MigrateUserLog()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Starting server on 127.0.0.1:%s ...", port)
	if err := r.SetupRouter().Run(":" + port); err != nil {
		log.Fatal(err)
	}
}
