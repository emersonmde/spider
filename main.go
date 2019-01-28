package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"

	"golang.org/x/net/html"
)

func main() {
	url := "https://google.com"

	links := getUrls(url)

	for link := range links {
		fmt.Printf("link = %+v\n", link)
	}
}

func getUrls(base string) <-chan string {
	baseUrl, err := url.Parse(base)
	if err != nil {
		log.Fatalf("Error parsing base url: %s", err)
	}
	out := make(chan string)
	page := getRequest(baseUrl)
	defer page.Close()
	links := parse(page)
	found := make(map[string]bool)
	go func() {
		for _, link := range links {
			//fmt.Printf("%v %v\n", baseUrl, link)
			if strings.Contains(link.Hostname(), baseUrl.Hostname()) && !found[link.RequestURI()] {
				found[link.RequestURI()] = true
				out <- link.String()
			}
		}
		close(out)
	}()
	return out
}

func getRequest(url *url.URL) io.ReadCloser {
	resp, err := http.Get(url.String())
	if err != nil {
		log.Fatal(err)
	}
	return resp.Body
}

func tokenUrl(t html.Token) (*url.URL, bool) {
	for _, a := range t.Attr {
		if a.Key == "href" {
			link, err := url.Parse(a.Val)
			if err != nil || !link.IsAbs() {
				return link, false
			}
			return link, true
		}
	}
	return nil, false
}

func parse(b io.Reader) []*url.URL {
	tokens := html.NewTokenizer(b)
	var urls []*url.URL
	for {
		tokenType := tokens.Next()
		switch {
		case tokenType == html.ErrorToken:
			return urls
		case tokenType == html.StartTagToken:
			t := tokens.Token()
			if !(t.Data == "a") {
				continue
			}
			url, ok := tokenUrl(t)
			if !ok {
				continue
			}
			urls = append(urls, url)
		}
	}
}
