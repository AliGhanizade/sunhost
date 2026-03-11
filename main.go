package main

import (
	"fmt"
	"sunhost/config"
	r "sunhost/router"
)

func main()  {

	config.InitDB()
	config.MigrateUser()
	config.MigrateSystemStat()
	config.MigrateUserLog()



	fmt.Println("Starting server on port 127.0.0.1:8080...")
	r.SetupRouter().Run(":8080")
}