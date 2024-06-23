package main

import "snake_ai/internal/application"

func main() {
	application.ParseFlags()
	if err := application.Run(); err != nil {
		panic(err)
	}
}
