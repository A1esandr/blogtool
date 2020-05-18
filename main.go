package main

import (
	appack "blogtool/src/app"
	"flag"
	"log"
	"os"
)

func main() {
	var urlFlag = flag.String("url", "", "Url for site. Like, http://example.com")
	url := os.Getenv("URL")
	if len(url) == 0 {
		url = *urlFlag
	}
	if len(url) == 0 {
		log.Fatalf("URL not set!")
	}
	backup := os.Getenv("MAKE_BACKUP") == "1"
	backupPath := os.Getenv("BACKUP_PATH")
	printAll := os.Getenv("PRINT_ALL") == "1"

	config := appack.Config{
		Url:        url,
		Backup:     backup,
		BackupPath: backupPath,
		PrintAll:   printAll,
	}

	app := appack.NewApp(config)
	app.Start()
	app.Print()
}
