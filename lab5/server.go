package main

import (
	"database/sql"
	"fmt"
	"github.com/SlyMarbo/rss"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"strings"
	"time"
)

var database *sql.DB

type news struct {
	title       string
	description string
}

func insertIntoDb(new news) {
	db, err := sql.Open("mysql", "iu9networkslabs:Je2dTYr6@tcp(students.yss.su:3306)/iu9networkslabs")
	if err != nil {
		log.Println(err)
	}
	database = db
	defer database.Close()
	query := database.QueryRow("SELECT EXISTS(SELECT `id` FROM `iu9Tarakanov` WHERE `title`=? OR `content`=?);", new.title, new.description)
	var isExists bool

	query.Scan(&isExists)

	if !isExists {
		database.Exec("INSERT INTO `iu9Tarakanov` (`title`, `content`) VALUES (?, ?);", new.title, new.description)
	} else {
		database.Exec("UPDATE `iu9Tarakanov` SET `content`=? WHERE `title`=?;", new.description, new.title)
		database.Exec("UPDATE `iu9Tarakanov` SET `title`=? WHERE `content`=?;", new.title, new.description)
	}

}

func rssparser() {
	rssObject, err := rss.Fetch("https://news.rambler.ru/rss/Namibia/")
	if err == nil {
		fmt.Printf("Title           : %s\n", rssObject.Title)
		fmt.Printf("Description     : %s\n", rssObject.Description)
		fmt.Printf("Link            : %s\n", rssObject.Link)
		fmt.Printf("Number of Items : %d\n", len(rssObject.Items))
		for v := range rssObject.Items {
			item := rssObject.Items[v]
			new_news := news{}
			new_news.title = strings.ReplaceAll(item.Title, "\u00A0", " ")
			new_news.description = strings.ReplaceAll(item.Summary, "\u00A0", " ")
			insertIntoDb(new_news)
			fmt.Println()
			fmt.Printf("Item Number : %d\n", v)
			fmt.Printf("Title       : %s\n", item.Title)
			fmt.Printf("Description : %s\n", item.Summary)

		}
	} else {
		fmt.Println(err)
	}
}

func main() {
	for {
		rssparser()
		time.Sleep(time.Second * 5)
	}

}
