package main

import (
	"flag"
	"os"

	"vyshkat/internal/app"
)

func main() {
	demo := flag.Bool("demo", false, "показать демо и выйти")
	flag.Parse()

	console := app.NewConsoleFromContainer(app.BuildContainer(), os.Stdin, os.Stdout)
	if *demo {
		console.RunDemo()
		return
	}
	console.Run()
}
