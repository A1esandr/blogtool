package app

import (
	"log"
	"os"
	"strings"
	"time"
)

type pathConfig struct {
}

type PathConfigurer interface {
	Configure(base, url string) string
}

func NewPathConfigurer() PathConfigurer {
	return &pathConfig{}
}

func (p *pathConfig) Configure(base, url string) string {
	backupPath := base
	if len(backupPath) > 0 && !strings.HasSuffix(backupPath, "/") {
		backupPath += "/"
	}
	backupPath += strings.Split(url, "/")[2]
	backupPath += "/"
	t := time.Now()
	backupPath += t.Format("2006-01-02")
	backupPath += "/"

	if _, err := os.Stat(backupPath); os.IsNotExist(err) {
		err = os.MkdirAll(backupPath, os.ModePerm)
		if err != nil {
			log.Fatalf(err.Error())
		}
	}

	return backupPath
}
