package main

import (
	"fmt"
	"github.com/SlyMarbo/rss"
)

func main() {
	rssObject, err := rss.Fetch("https://briansk.ru/rss20_briansk.xml")
	if err == nil {
		fmt.Printf("Title           : %s\n", rssObject.Title)
		fmt.Printf("Description     : %s\n", rssObject.Description)
		fmt.Printf("Link            : %s\n", rssObject.Link)
		fmt.Printf("Number of Items : %d\n", len(rssObject.Items))
		for v := range rssObject.Items {
			item := rssObject.Items[v]
			fmt.Println()
			fmt.Printf("Item Number : %d\n", v)
			fmt.Printf("Title       : %s\n", item.Title)
			fmt.Printf("Link        : %s\n", item.Link)
			fmt.Printf("Url picture : %s\n", item.Enclosures[0].URL)
			fmt.Printf("Description : %s\n", item.Summary)
		}
	} else {
		fmt.Println(err)
	}
}
