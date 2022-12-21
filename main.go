package main

import (
	"fmt"
	"os"

	"camis-v2-api/config"
	"camis-v2-api/routes"
)

func main() {
	fmt.Println("Hello - CAMIS")
	config.InitializeConfig()
	config.ConnectDb()
	defer config.SESSION.Close()
	defer config.DB.Close()
	// controller.SocketConnection()

	server := routes.InitRoutes()
	// ctx := context.Background()
	// go controller.StatusCheckLoop(ctx)
	server.Run(":" + os.Getenv("server_port"))
}
