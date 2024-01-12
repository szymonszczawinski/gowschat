package main

import (
	"gowschat/server"
	"log"

	"github.com/joho/godotenv"
)

func init() {
	godotenv.Load()
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}

func main() {
	log.Println("Hello Chat Server")
	server.Run()
}
