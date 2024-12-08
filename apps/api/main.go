package main

import (
	"apps/api/handler"
)

func main() {
	server := handler.NewServer()
	server.Start()
}
