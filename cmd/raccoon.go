package main

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/snwzt/raccoon/internal/manager"
)

func main() {
	customlogger := zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout, NoColor: true}).With().Timestamp().Logger()

	manager.Execute(
		&customlogger,
		os.Exit,
		os.Args[1:],
	)
}
