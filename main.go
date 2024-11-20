package main

import (
	"GRIMIDIE/cli"
	"GRIMIDIE/config"
	"GRIMIDIE/handler"

	_ "github.com/lib/pq"
)

func main() {
	db := config.ConnectDB()

	handler := handler.NewHandler(db)

	cli := cli.NewCLI(handler)
	cli.Init()
}
