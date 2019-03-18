package main

import (
	"log"
	"os"

	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "debased-cli"
	app.Usage = "the debased command line interface"

	// Run the app
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}

}
