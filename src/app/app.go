package app

import (
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"golang.org/x/net/html"
)

type App interface {
	Start()
	Print()
}

type app struct {
	url        string
	config     Config
	backupPath string
	links      []string
	rawLinks   []string
	loaded     []string
	left       []string
	lock       sync.Mutex
}

func NewApp(config Config) App {
	backupPath := config.BackupPath
	if len(backupPath) > 0 && !strings.HasSuffix(backupPath, "/") {
		backupPath += "/"
	}
	backupPath += strings.Split(config.Url, "/")[2]
	backupPath += "/"
	t := time.Now()
	backupPath += t.Format("2006-01-02")
	backupPath += "/"

	if _, err := os.Stat(backupPath); os.IsNotExist(err) && config.Backup {
		err = os.MkdirAll(backupPath, os.ModePerm)
		if err != nil {
			log.Fatalf(err.Error())
		}
	}

	return &app{
		url:        config.Url,
		config:     config,
		backupPath: backupPath,
	}
}

func (a *app) Start() {
	a.process(a.url)
	a.next()
}

func (a *app) process(url string) {
	var body []byte
	for {
		b, err := a.load(url)
		if err != nil {
			log.Println(err.Error())
		}
		if len(b) > 1000 || !strings.HasSuffix(url, ".html") {
			body = b
			break
		}
		log.Println("Error loading", url)
		time.Sleep(time.Duration(300+rand.Intn(1000)) * time.Millisecond)
	}
	if a.config.Backup && strings.HasSuffix(url, ".html") {
		a.backup(body, url)
	}
	doc, err := html.Parse(strings.NewReader(string(body)))
	if err != nil {
		log.Fatal(err)
	}
	a.parse(doc)
}

func (a *app) backup(file []byte, url string) {
	name := strings.ReplaceAll(url, "://", "")
	name = strings.ReplaceAll(name, "/", "_")
	err := ioutil.WriteFile(a.backupPath+name, file, 0644)
	if err != nil {
		panic(err)
	}
}

func (a *app) Print() {
	fmt.Println("Print result")
	for _, link := range a.links {
		fmt.Println(link)
	}
	fmt.Println("Total found", len(a.rawLinks), "raw links")
	fmt.Println("Total found", len(a.links), "links")
}

func (a *app) load(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer func() {
		err = resp.Body.Close()
		if err != nil {
			log.Println(err.Error())
		}
	}()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	a.lock.Lock()
	a.loaded = append(a.loaded, url)
	a.lock.Unlock()
	log.Println("Loaded", url)
	return body, nil
}

func (a *app) parse(n *html.Node) {
	if n.Type == html.ElementNode && n.Data == "a" {
		for _, at := range n.Attr {
			if at.Key == "href" {
				a.lock.Lock()
				if strings.HasPrefix(at.Val, a.url) && !a.exist(at.Val) {
					a.add(at.Val)
				}
				a.lock.Unlock()
				break
			}
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		a.parse(c)
	}
}

func (a *app) exist(url string) bool {
	for _, link := range a.rawLinks {
		if url == link {
			return true
		}
	}
	return false
}

func (a *app) processed(url string) bool {
	for _, link := range a.loaded {
		if url == link {
			return true
		}
	}
	return false
}

func (a *app) add(url string) {
	a.rawLinks = append(a.rawLinks, url)
	if strings.HasSuffix(url, ".html") {
		a.links = append(a.links, url)
	}
	if !a.processed(url) {
		a.left = append(a.left, url)
	}
}

func (a *app) next() {
	if len(a.left) == 0 {
		return
	}

	a.lock.Lock()

	log.Println("Processed", len(a.loaded))
	log.Println("Left", len(a.left))

	var wg sync.WaitGroup

	for _, link := range a.left {
		wg.Add(1)
		go func(url string) {
			a.process(url)
			wg.Done()
		}(link)
	}
	a.left = []string{}
	a.lock.Unlock()
	wg.Wait()
	a.next()
}
