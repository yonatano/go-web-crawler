package main

import (
	"fmt"
	// "regexp"
	// "io"
	"bytes"
	"github.com/PuerkitoBio/purell"
	"github.com/oleiade/lane"
	"gopkg.in/xmlpath.v2"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

type Crawler struct {
	queue *lane.Queue
}

type Page struct {
	url    string
	source []byte
}

func (c *Crawler) Crawl() {
	var linkQueue *lane.Queue = c.queue
	fmt.Println(linkQueue)
	for {
		link := linkQueue.Dequeue()
		if linkstr, ok := link.(string); ok {
			//fetch data with link
			p := Page{url: linkstr}
			p.Fetch()
			fetched := p.GetLinks()
			//add new links to queue
			for _, l := range fetched {
				linkQueue.Enqueue(l)
			}
			//report
			fmt.Println("crawled: %s added %i new links", linkstr, len(fetched))
		} else {
			break
		}
	}
}

func (p *Page) Fetch() {
	//make request
	resp, err := http.Get(p.url)
	if err != nil {
		// fail and requeue for later
		panic(err)
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
		link := iter.Node().String()
		//fix and normalize url
		link = FixUrl(link, p.url)
		links = append(links, link)
	}
	return links
}

func FixUrl(link string, parent string) string {
	//normalize url
	normalized := purell.MustNormalizeURLString(link, purell.FlagsSafe)
	u, err := url.Parse(normalized)
	if err != nil {
		log.Fatal(err)
	}
	base, err := url.Parse(parent)
	if err != nil {
		log.Fatal(err)
	}
	return base.ResolveReference(u).String()
}

func main() {
	var seed *lane.Queue = lane.NewQueue()
	seed.Enqueue("http://wikipedia.com/")
	var crawler = Crawler{queue: seed}
	crawler.Crawl()
	// re := regexp.MustCompile("href=""")
	// fmt.Println(re.MatchString("paranormal"))
}
