package main

import (
	"context"
	"log"
	"os"

	"github.com/atian25/projj-go/cmd"
)

func main() {
	app := cmd.NewApp()
	
	if err := app.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}