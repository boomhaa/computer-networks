package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/websocket"
)

type Article struct {
	ID      int
	Title   string
	Content string
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func dbConn() (db *sql.DB, err error) {
	dsn := "iu9networkslabs:Je2dTYr6@tcp(students.yss.su:3306)/iu9networkslabs"
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil
}

func fetchArticles() ([]Article, error) {
	db, err := dbConn()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	rows, err := db.Query("SELECT id, title, content FROM iu9Tarakanov")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var articles []Article
	for rows.Next() {
		var article Article
		if err := rows.Scan(&article.ID, &article.Title, &article.Content); err != nil {
			return nil, err
		}
		articles = append(articles, article)
	}

	return articles, nil
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "Failed to upgrade connection", http.StatusInternalServerError)
		return
	}
	defer conn.Close()

	for {
		articles, err := fetchArticles()
		if err != nil {
			log.Println("Error fetching articles:", err)
			return
		}

		err = conn.WriteJSON(articles)
		if err != nil {
			log.Println("Error writing to websocket:", err)
			return
		}

		time.Sleep(5 * time.Second)
	}
}

func main() {
	http.HandleFunc("/ws", wsHandler)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := template.New("index").Parse(`
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Articles</title>
    <script>
        let socket = new WebSocket("ws://" + window.location.host + "/ws");
        socket.onmessage = function(event) {
            let articles = JSON.parse(event.data);
            let container = document.getElementById("articles");
            container.innerHTML = "";
            articles.forEach(article => {
                let articleDiv = document.createElement("div");
                let title = document.createElement("h2");
                title.textContent = article.Title;
                let content = document.createElement("p");
                content.textContent = article.Content;
                articleDiv.appendChild(title);
                articleDiv.appendChild(content);
				articleDiv.appendChild(document.createElement("hr"));
                container.appendChild(articleDiv);
            });
        };
    </script>
</head>
<body>
    <h1>Articles</h1>
    <div id="articles"></div>
</body>
</html>
`)
		if err != nil {
			http.Error(w, "Failed to create template", http.StatusInternalServerError)
			return
		}
		tmpl.Execute(w, nil)
	})

	fmt.Println("Server is running on 9786")
	log.Fatal(http.ListenAndServe(":9786", nil))
}
