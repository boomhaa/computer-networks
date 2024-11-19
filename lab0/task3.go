package main

import (
	"fmt"      
	"log"      
	"net/http" 
	"strings"  
	"github.com/SlyMarbo/rss"
)

func HomeRouterHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()      
	fmt.Println(r.Form) 
	fmt.Println("path", r.URL.Path)
	fmt.Println("scheme", r.URL.Scheme)
	fmt.Println(r.Form["url_long"])
	for k, v := range r.Form {
		fmt.Println("key:", k)
		fmt.Println("val:", strings.Join(v, ""))
	}

	for {
		rssObject, err := rss.Fetch("https://briansk.ru/rss20_briansk.xml")
		if err == nil {
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			fmt.Fprintf(w, "<!DOCTYPE html><html><head><title>RSS Feed</title></head><body>")
			fmt.Fprintf(w, "<h1>Description:</h1><p><strong>%s</strong></p>", rssObject.Description)
			fmt.Fprintf(w, "<h1>Number of Items:</h1><p><strong>%d</strong></p>", len(rssObject.Items))
			for v := range rssObject.Items {
				item := rssObject.Items[v]
				fmt.Fprintf(w, "<h1>%s</h1>", item.Title)
				fmt.Fprintf(w, "<p>%s</p>", item.Summary)
				fmt.Fprint(w, "<hr>")

			}

		}
	}

}

func main() {
	http.HandleFunc("185.104.251.226/", HomeRouterHandler) // установим роутер
	err := http.ListenAndServe(":9999", nil)               // задаем слушать порт
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
