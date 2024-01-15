package main

import (
	"context"
	"fmt"
	"gowschat/server"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
)

func init() {
	godotenv.Load()
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}

func main() {
	log.Println("Hello Chat Server")
	server.RunApp()
}

func main2() {
	c1 := context.Background()

	c2, cancel2 := context.WithCancel(c1)

	go functionC2(c2)
	quit := make(chan os.Signal, 1)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown...")
	cancel2()
	time.Sleep(1 * time.Second)
	fmt.Println("Done")
	defer func() {
		log.Println("exit main")
	}()
}

func functionC2(ctx context.Context) {
	log.Println("run function C2")
	defer func() {
		log.Println("exiting function C2")
	}()
	ticker := time.NewTicker(1 * time.Second)
	for {
		select {
		case <-ticker.C:
			log.Println("Tick C2")
		case <-ctx.Done():
			log.Println("kill function C2")
			return

		}
	}
}
