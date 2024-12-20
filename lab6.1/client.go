package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/jlaffaye/ftp"
	_ "github.com/jlaffaye/ftp"
	"log"
	"net/http"
	"strings"
	"time"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}
var c *ftp.ServerConn
var res string
var curDir = "./"

func main() {

	http.HandleFunc("/work", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./public/work/index.html")
	})
	http.Handle("/", http.FileServer(http.Dir("./public/logined")))
	http.HandleFunc("/login", login)
	http.HandleFunc("/command", handleCommand)
	http.HandleFunc("/ws", handleWebSocket)
	log.Fatal(http.ListenAndServe(":8021", nil))

}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Ошибка при апгрейде WebSocket:", err)
		return
	}
	defer conn.Close()
	for {
		if err := conn.WriteJSON(res); err != nil {
			log.Println("Ошибка при отправке JSON через WebSocket:", err)
		}
		time.Sleep(time.Second)
	}

}

func handleCommand(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	type commandReq struct {
		Command string `json:"command"`
	}
	var req commandReq
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	cmd := strings.Split(req.Command, " ")
	command := strings.Split(req.Command, " ")[0]
	path := ""
	if len(cmd) == 2 {
		path = strings.Split(req.Command, " ")[1]
	}

	log.Println("Получена команда:", req.Command)
	if command == "LS" {
		res = ""
		list, err := c.List("./")
		if err != nil {
			res = err.Error()
			fmt.Println(err)
			return
		}
		for _, v := range list {
			res += v.Name + " " + v.Type.String() + "\n"

		}
	} else if command == "MKDIR" {
		err = c.MakeDir(path)
		if err != nil {
			res = err.Error()
			fmt.Println(err)
			return
		}
		res = "Directory with name " + path + " was created successfully\n"
	} else if command == "CD" {
		curDir += path + "/"
		err = c.ChangeDir(path)
		if path == ".." {
			new_path := strings.Split(curDir, "/")
			curDir = strings.Join(new_path[:len(new_path)-3], "/")
			curDir += "/"
		}

		fmt.Println(curDir)
		if err != nil {
			res = err.Error()
			fmt.Println(err)
			return
		}
		res = "Current directory is " + curDir + "\n"
	} else if command == "RMDIR" {
		err = c.RemoveDir(path)
		if err != nil {
			res = err.Error()
			fmt.Println(err)
			return
		}
		res = "Directory with name " + path + " was removed successfully\n"
	} else if command == "RM" {
		err = c.Delete(path)
		if err != nil {
			res = err.Error()
			fmt.Println(err)
			return
		}
		res = "File with name " + path + " was removed successfully\n"
	} else if command == "RMREC" {
		err = c.RemoveDirRecur(path)
		if err != nil {
			res = err.Error()
			fmt.Println(err)
			return
		}
		res = "Directory with name " + path + " was recursively removed recusive successfully\n"
	} else {
		res = "Unknown command: " + command
	}
	fmt.Println(command + " done")
	w.WriteHeader(http.StatusOK)
}

func login(w http.ResponseWriter, r *http.Request) {
	log.Println("попытка логина")
	r.Header.Add("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method != http.MethodPost {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	type data struct {
		Host     string `json:"host"`
		Login    string `json:"login"`
		Password string `json:"password"`
	}

	var newData data
	err := json.NewDecoder(r.Body).Decode(&newData)
	if err != nil {
		http.Error(w, "Неверный формат запроса", http.StatusBadRequest)
		return
	}
	fmt.Println(newData.Host, newData.Login, newData.Password)

	err = conFtp(newData.Host, newData.Login, newData.Password)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	res = ""
	log.Println("Успешный вход")
	w.WriteHeader(http.StatusOK)
}
func conFtp(host, login, password string) (err error) {
	c, err = ftp.Dial(host, ftp.DialWithTimeout(5*time.Second))
	if err != nil {
		log.Println(err)
		return err
	}

	err = c.Login(login, password)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}
