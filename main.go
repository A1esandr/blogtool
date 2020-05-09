package main

import (
	app2 "blogtool/src/app"
	"log"
	"os"
)

func main() {
	url := os.Getenv("URL")
	if len(url) == 0 {
		log.Fatalf("URL not set!")
	}

	app := app2.NewApp(url)
	app.Parse()
	app.Print()
}
