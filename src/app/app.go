package app

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"golang.org/x/net/html"
)

type App interface {
	Parse()
	Print()
}

type app struct {
	url      string
	links    []string
	rawLinks []string
	loaded   []string
}

func NewApp(url string) App {
	return &app{
		url: url,
	}
}

func (a *app) Parse() {
	body, err := a.load(a.url)
	if err != nil {
		log.Println(err.Error())
	}
	doc, err := html.Parse(strings.NewReader(string(body)))
	if err != nil {
		log.Fatal(err)
	}
	a.parse(doc)
}

func (a *app) Print() {
	for _, link := range a.rawLinks {
		fmt.Println(link)
	}
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
	return body, nil
}

func (a *app) parse(n *html.Node) {
	if n.Type == html.ElementNode && n.Data == "a" {
		for _, at := range n.Attr {
			if at.Key == "href" {
				if !a.exist(at.Val) {
					a.add(at.Val)
				}
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

func (a *app) add(url string) {
	a.rawLinks = append(a.rawLinks, url)
}
