package main

import (
	"fmt"
	// "regexp"
	// "io"
	"bytes"
	"gopkg.in/xmlpath.v2"
	"io/ioutil"
	"net/http"
)

type Crawler struct {
	queue []string
}

// func (c *Crawler)Enqueue

type Page struct {
	url    string
	source []byte
}

func (p *Page) Fetch() {
	//make request
	resp, err := http.Get(p.url)
	if err != nil {
		// handle error
		fmt.Println("error: ", err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	p.source = body
}

func (p *Page) GetLinks() []string {
	source := bytes.NewReader(p.source)
	node, err := xmlpath.ParseHTML(source)
	if err != nil {
		panic(err)
	}
	//extract links
	path := xmlpath.MustCompile("//a/@href")
	iter := path.Iter(node)

	var links []string
	for iter.Next() {
		url := iter.Node().String()
		links = append(links, url)
	}
	return links
}

// func (p *page)

func main() {
	p := Page{url: "http://wikipedia.org/"}
	p.Fetch()
	links := p.GetLinks()
	fmt.Println(links)
	//fmt.Println("source:\n", p.source)
	// re := regexp.MustCompile("href=""")
	// fmt.Println(re.MatchString("paranormal"))
}
