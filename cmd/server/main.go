package main

import "snakeai/internal/application"

func main() {
	application.ParseFlags()
	if err := application.Run(); err != nil {
		panic(err)
	}
}
