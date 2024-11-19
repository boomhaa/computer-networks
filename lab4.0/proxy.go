package main

import (
	"bytes"
	"fmt"
	"github.com/hashicorp/go-getter/v2"
	"golang.org/x/net/context"
	"golang.org/x/net/html"
	"net/http"
	"strings"
)

var ip = "192.168.31.116:9457"
var global_ip = "http://46.138.248.69:9457"

func replace(doc *html.Node, baseURL string) error {
	var clawler func(*html.Node)
	clawler = func(n *html.Node) {
		if n.Type == html.ElementNode {
			for i, attr := range n.Attr {
				if attr.Key == "href" || attr.Key == "src" || attr.Key == "action" {
					if attr.Val != "/" {
						if !strings.HasPrefix(attr.Val, "http") {
							fmt.Println("Поменял", n.Attr[i].Val, "на", global_ip+"/"+baseURL+attr.Val)
							n.Attr[i].Val = global_ip + "/" + baseURL + attr.Val
						} else {
							fmt.Println("Поменял", n.Attr[i].Val, "на", global_ip+"/"+attr.Val)
							n.Attr[i].Val = global_ip + "/" + attr.Val
						}
					} else {
						if !strings.HasPrefix(attr.Val, "http") {
							fmt.Println("Поменял", n.Attr[i].Val, "на", global_ip+"/"+baseURL)
							n.Attr[i].Val = global_ip + "/" + baseURL + attr.Val
						} else {
							fmt.Println("Поменял", n.Attr[i].Val, "на", global_ip+"/"+n.Attr[i].Val)
							n.Attr[i].Val = global_ip + "/" + attr.Val
						}
					}

				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			clawler(c)
		}
	}
	clawler(doc)
	return nil

}

func handler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	fmt.Println(path)
	path = strings.TrimPrefix(path, "/")
	path = strings.Replace(path, "/", "//", 1)
	isHtml := false
	if strings.HasSuffix(path, ".html") || (!strings.HasSuffix(path, ".css") &&
		!strings.HasSuffix(path, ".js") && !strings.HasSuffix(path, ".png") &&
		!strings.HasSuffix(path, ".gif") && !strings.HasSuffix(path, ".jpg") && !strings.HasSuffix(path, ".jpeg") &&
		!strings.HasSuffix(path, ".ico") && !strings.HasSuffix(path, ".svg") && !strings.HasSuffix(path, ".ttf") && !strings.HasSuffix(path, ".woff") && !strings.HasSuffix(path, ".woff2")) {

		isHtml = true
	}
	if isHtml {
		get, err := http.Get(path)
		if err != nil {
			fmt.Println(err)
			return
		}
		doc, err := html.Parse(get.Body)
		fmt.Println("Получил страницу")

		baseURL := path[:strings.LastIndex(path, "/")+1]
		err = replace(doc, baseURL)
		fmt.Println("Поменял ссылки")
		if err != nil {
			fmt.Println(err)
			return
		}
		var buf bytes.Buffer
		if err := html.Render(&buf, doc); err != nil {
			http.Error(w, "Ошибка генерации HTML: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.Write(buf.Bytes())
	} else {
		fmt.Println("Картинка")
		filename := strings.Split(path, "/")
		_, err := getter.GetFile(context.Background(), "./static/"+filename[len(filename)-1], path)
		if err != nil {
			fmt.Println(err)
		} else {
			http.ServeFile(w, r, "./static/"+filename[len(filename)-1])
		}
	}
}

func main() {

	http.HandleFunc("/", handler)

	fmt.Println("Прокси-сервер запущен на порту 9457")
	err := http.ListenAndServe(ip, nil)
	if err != nil {
		fmt.Println("Ошибка запуска сервера:", err)
		return
	}
}
