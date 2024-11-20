// Go & MySQL CRUD Example
package main

import (
	// "go-mysql/cli"
	// "go-mysql/config"
	// "go-mysql/handler"

	"GRIMIDIE/config"

	_ "github.com/lib/pq"
)

func main() {
	// db := 
	config.ConnectDB()

	// handler := handler.NewHandler(db)

	// cli := cli.NewCLI(handler)
	// cli.Init()
}
