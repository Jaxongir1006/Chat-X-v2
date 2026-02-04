package main

import (
	"log"
	"os"

	"github.com/Jaxongir1006/Chat-X-v2/internal/app"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Please provide a command: http | superuser")
	}

	cmd := os.Args[1]

	app.Run(cmd)
}
