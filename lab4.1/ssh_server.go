package main

import (
	"fmt"
	"github.com/gliderlabs/ssh"
	"io"
	"log"
	"os/exec"
	"strings"
	"sync"
)

var (
	users      = map[string]string{"root": "1234"}
	usersMutex sync.Mutex
)

func passwordAuth(ctx ssh.Context, password string) bool {
	usersMutex.Lock()
	defer usersMutex.Unlock()

	if pass, ok := users[ctx.User()]; ok && pass == password {
		return true
	}
	return false
}

func handleAddUser(args []string, session ssh.Session) {
	if session.User() != "root" {
		io.WriteString(session, "Ошибка: добавление пользователя доступно только для root\n")
		return
	}
	if len(args) < 2 {
		io.WriteString(session, "Ошибка: требуется указать имя пользователя и пароль\n")
		return
	}

	username := args[0]
	password := args[1]

	usersMutex.Lock()
	defer usersMutex.Unlock()

	if _, exists := users[username]; exists {
		io.WriteString(session, "Ошибка: пользователь уже существует\n")
	} else {
		users[username] = password
		io.WriteString(session, fmt.Sprintf("Пользователь %s успешно добавлен\n", username))
	}
}

func main() {
	server := ssh.Server{
		Addr:            ":9456",
		PasswordHandler: passwordAuth,
	}
	server.Handler = func(s ssh.Session) {
		command := s.RawCommand()
		args := strings.Split(command, " ")

		if len(args) > 0 && args[0] == "adduser" {
			handleAddUser(args[1:], s)
		} else {
			cmd := exec.Command(args[0], args[1:]...)
			stdout, err := cmd.Output()
			if err != nil {
				io.WriteString(s, fmt.Sprintf("ERROR OCCURED: %s\n", err.Error()))
			} else {
				io.WriteString(s, fmt.Sprintf("%s\n", string(stdout)))
			}
		}
	}

	log.Println("Запуск SSH сервера на порту 9456...")
	log.Fatal(server.ListenAndServe())
}
