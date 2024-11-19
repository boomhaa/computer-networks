package main

import (
	"net/http"
	"strconv"

	log "github.com/mgutz/logxi/v1"
	"golang.org/x/net/html"
)

func getAttr(node *html.Node, key string) string {
	for _, attr := range node.Attr {
		if attr.Key == key {
			return attr.Val
		}
	}
	return ""
}

func getChildren(node *html.Node) []*html.Node {
	var children []*html.Node
	for c := node.FirstChild; c != nil; c = c.NextSibling {
		children = append(children, c)
	}
	return children
}

func isElem(node *html.Node, tag string) bool {
	return node != nil && node.Type == html.ElementNode && node.Data == tag
}

func isText(node *html.Node) bool {
	return node != nil && node.Type == html.TextNode
}

func isDiv(node *html.Node, class string) bool {
	return isElem(node, "div") && getAttr(node, "class") == class
}

type Item struct {
	Ref, Title string
}

func readItem(item *html.Node) *Item {
	if a := item.FirstChild; isElem(a, "a") {
		return &Item{
			Ref:   getAttr(a, "href"),
			Title: getAttr(a, "aria-label"),
		}
	}
	return nil
}

func search(node *html.Node) []*Item {
	log.Info("begin of parsing")
	if isElem(node, "div") && getAttr(node, "id") == "app" {
		log.Info("нашел app")
		var items []*Item
		for c := node.FirstChild; c != nil; c = c.NextSibling {
			if isDiv(c, "dTSkA6xB commercial-branding") {
				log.Info("нашел dTSkA6xB commercial-branding")
				for d := c.FirstChild; d != nil; d = d.NextSibling {
					if isDiv(d, "AuRBdDZg") {
						log.Info("нашел AuRBdDZg")
						for k := d.FirstChild; k != nil; k = k.NextSibling {
							if isElem(k, "section") {
								log.Info("нашел section")
								for e := k.FirstChild; e != nil; e = e.NextSibling {
									if isDiv(e, "cGZPyk4_") {
										log.Info("нашел cGZPyk4_")
										for g := e.FirstChild; g != nil; g = g.NextSibling {
											if isDiv(g, "zT5wwAPN fQtJ19Ei") {
												log.Info("нашел zT5wwAPN fQtJ19Ei")
												for y := g.FirstChild; y != nil; y = y.NextSibling {
													if isDiv(y, "XSvLK2D0 abGoxuyb") {
														log.Info("нашел XSvLK2D0 abGoxuyb")
														items = append(items, readItem(y))
													}

												}
												log.Info("total: " + strconv.Itoa(len(items)) + " news")
												return items
											}
										}
									}
								}
							}
						}
					}
				}

			}
		}
	}
	for c := node.FirstChild; c != nil; c = c.NextSibling {
		if items := search(c); items != nil {
			return items
		}
	}
	return nil
}

func downloadNews() []*Item {
	log.Info("sending request to https://news.rambler.ru/latest/")
	if response, err := http.Get("https://news.rambler.ru/latest/"); err != nil {
		log.Error("request to https://news.rambler.ru/latest/ failed", "error", err)
	} else {
		defer response.Body.Close()
		status := response.StatusCode
		log.Info("got response from https://news.rambler.ru/latest/", "status", status)
		if status == http.StatusOK {
			if doc, err := html.Parse(response.Body); err != nil {
				log.Error("invalid HTML from https://news.rambler.ru/latest/", "error", err)
			} else {
				log.Info("HTML from https://news.rambler.ru/latest/ parsed successfully")
				return search(doc)
			}
		}
	}
	return nil
}
