package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

const xkcdRandomURL = "https://c.xkcd.com/random/comic/"
const xkcdComicSelector = "#comic img"

type comic struct {
	title, alt, src string
}

func getComicFromSelection(s *goquery.Selection) (comic, error) {
	var c comic

	attrMap := map[string]string{"title": "", "alt": "", "src": ""}
	for key := range attrMap {
		val, found := s.Attr(key)
		if !found {
			return c, errors.New("html document error: missing comic attr: " + key)
		}
		attrMap[key] = val
	}

	c.title = attrMap["title"]
	c.alt = attrMap["alt"]
	c.src = strings.Replace(attrMap["src"], "//", "https://", 1)

	return c, nil
}

func main() {
	// Request the HTML
	res, err := http.Get(xkcdRandomURL)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	// Find the comic image
	comicImg := doc.Find(xkcdComicSelector)
	for i := range comicImg.Nodes {
		img := comicImg.Eq(i)
		comic, err := getComicFromSelection(img)
		if err != nil {
			log.Fatal(err)
		}

		// Save the image
		// TODO
		fmt.Printf("Title: %s\nAlt: %s\nSrc: %s\n", comic.title, comic.alt, comic.src)
	}

}
