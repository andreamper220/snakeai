package logger

import (
	"fmt"
	"os"

	"go.uber.org/zap"
)

var Log *zap.SugaredLogger

func Initialize() error {
	logger, err := zap.NewDevelopment()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	Log = logger.Sugar()

	return nil
}
