package main

import (
	server "github.com/vapusdata-oss/aistudio/aistudio/server"
)

func main() {
	// This is the main function to start the platform server
	server := server.GrpcServer()
	server.Run()
}
