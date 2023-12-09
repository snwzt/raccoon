package main

import (
	"os"

	"github.com/snwzt/raccoon/internal/manager"
)

func main() {
	manager.Execute(
		os.Exit,
		os.Args[1:],
	)
}
