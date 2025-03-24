package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	link "github.com/Kacperbrozyna/learning_go_repo/http_link_parser/link_parser"
)

const xmlns = "http://www.sitemaps.org/schemas/sitemap/0.9"

type loc struct {
	Value string `xml:"loc"`
}

type urlSet struct {
	Urls  []loc  `xml:"url"`
	Xmlns string `xml:"xmlns,attr"`
}

func main() {
	url_flag := flag.String("url", "https://gophercises.com", "The url that you want to build a site map for")
	max_depth := flag.Int("depth", 3, "The maximum number of links deep to traverse")
	flag.Parse()

	pages := bfs(*url_flag, *max_depth)
	to_xml := urlSet{
		Xmlns: xmlns,
	}

	for _, page := range pages {
		to_xml.Urls = append(to_xml.Urls, loc{page})
	}

	fmt.Print(xml.Header)
	encoder := xml.NewEncoder(os.Stdout)
	encoder.Indent("", "  ")
	if err := encoder.Encode(to_xml); err != nil {
		panic(err)
	}

	fmt.Println()
}

func bfs(url_string string, max_depth int) []string {

	seen := make(map[string]struct{})
	var queue map[string]struct{}
	next_queue := map[string]struct{}{
		url_string: {},
	}

	for i := 0; i <= max_depth; i++ {
		queue, next_queue = next_queue, make(map[string]struct{})
		if len(queue) == 0 {
			break
		}
		for url := range queue {
			if _, ok := seen[url]; ok {
				continue
			}

			seen[url] = struct{}{}

			for _, link := range get(url) {
				if _, ok := seen[link]; !ok {
					next_queue[link] = struct{}{}
				}
			}
		}
	}

	ret := make([]string, 0, len(seen))
	for url := range seen {
		ret = append(ret, url)
	}

	return ret
}

func get(url_string string) []string {
	response, err := http.Get(url_string)
	if err != nil {
		return []string{}
	}
	defer response.Body.Close()

	request_url := response.Request.URL
	base_url := &url.URL{
		Scheme: request_url.Scheme,
		Host:   request_url.Host,
	}

	base := base_url.String()

	return filter(hrefs(response.Body, base), withPrefix(base))
}

func hrefs(reader io.Reader, base string) []string {
	links, err := link.Parse(reader)
	if err != nil {
		return []string{}
	}
	var ret []string
	for _, link := range links {
		switch {
		case strings.HasPrefix(link.Href, "/"):
			ret = append(ret, base+link.Href)
		case strings.HasPrefix(link.Href, "http"):
			ret = append(ret, link.Href)
		}
	}

	return ret
}

func filter(links []string, keep_fn func(string) bool) []string {
	var ret []string

	for _, link := range links {

		if keep_fn(link) {
			ret = append(ret, link)
		}
	}
	return ret
}

func withPrefix(prefix string) func(string) bool {
	return func(link string) bool {
		return strings.HasPrefix(link, prefix)
	}
}
