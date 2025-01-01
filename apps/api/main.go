package main

import (
	"apps/api/server"
)

func main() {
	server := server.NewServer()
	server.Start()
}
