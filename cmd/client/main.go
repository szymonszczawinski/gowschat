package main

import (
	"log"

	"github.com/joho/godotenv"
)

func init() {
	godotenv.Load()
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}

func main() {
	log.Println("Hello Clinet")
}
