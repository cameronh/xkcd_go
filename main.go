package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"golang.org/x/net/html"
)

// Image ... crawls the html response for the meta tag containing the source image
func Image(doc *html.Node) (*html.Node, error) {
	var image *html.Node
	var crawler func(*html.Node)
	crawler = func(node *html.Node) {
		if node.Type == html.ElementNode && node.Data == "meta" {

			image = node
			return
		}
		for child := node.FirstChild; child != nil; child = child.NextSibling {
			crawler(child)
		}
	}
	crawler(doc)
	if image != nil {
		return image, nil
	}
	return nil, errors.New("Missing <meta> in the node tree")
}

func main() {
	res, err := http.Get("https://c.xkcd.com/random/comic/")
	if err != nil {
		log.Fatal(err)
	}

	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}

	doc, _ := html.Parse(strings.NewReader(string(body)))
	image, err := Image(doc)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(image)
}
