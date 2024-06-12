package main

import "snake_ai/internal/server"

func main() {
	server.ParseFlags()
	if err := server.Run(); err != nil {
		panic(err)
	}
}
