package main

import (
	editorserver "github.com/andreamper220/snakeai/internal/editor/application"
)

func main() {
	if err := editorserver.Run(); err != nil {
		panic(err)
	}
}
