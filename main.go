package main

import (
	"fmt"
)

func main() {
	app := newFileServer()
	app.configureRoutes()

	go app.listenAndServe(":8080")
	fmt.Println("Server listening on Port 8080")

	fmt.Println("Server started, hit Enter-key to close")
	fmt.Scanln()
	fmt.Println("Shuting down...")
}
