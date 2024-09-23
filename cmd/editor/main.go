package main

import (
	editorserver "github.com/andreamper220/snakeai/internal/editor/application"
)

func main() {
	if err := editorserver.Run(50051, false); err != nil {
		panic(err)
	}
}
