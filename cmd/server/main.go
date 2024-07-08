package main

import "github.com/andreamper220/snakeai/internal/application"

func main() {
	application.ParseFlags()
	if err := application.Run(false); err != nil {
		panic(err)
	}
}
