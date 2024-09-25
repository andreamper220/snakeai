package main

import (
	gameserver "github.com/andreamper220/snakeai/internal/server/application"
)

func main() {
	gameserver.ParseFlags()
	if err := gameserver.Run(false); err != nil {
		panic(err)
	}
}
