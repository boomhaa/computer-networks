package main

import (
	"fmt"
	"golang.org/x/crypto/ssh"
	"log"
	"os"
	"strings"
)

const (
	sshUser       = "root"
	sshPassword   = "1234"
	sshServerHost = "185.104.251.226:9456"
)

func createSSHClient(user, password, host string) (*ssh.Client, error) {
	config := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	client, err := ssh.Dial("tcp", host, config)
	if err != nil {
		return nil, fmt.Errorf("не удалось подключиться к серверу: %v", err)
	}
	return client, nil
}

func runCommand(client *ssh.Client, command string) (string, error) {
	session, err := client.NewSession()
	if err != nil {
		return "", fmt.Errorf("не удалось создать сессию: %v", err)
	}
	defer session.Close()

	var outputBuf strings.Builder
	session.Stdout = &outputBuf
	session.Stderr = &outputBuf

	if err = session.Run(command); err != nil {
		return "", fmt.Errorf("ошибка при выполнении команды: %v", err)
	}

	return outputBuf.String(), nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Использование: go run ssh_client.go <команда>")
		os.Exit(1)
	}
	command := strings.Join(os.Args[1:], " ")

	client, err := createSSHClient(sshUser, sshPassword, sshServerHost)
	if err != nil {
		log.Fatalf("Ошибка подключения: %v", err)
	}
	defer client.Close()

	fmt.Printf("Выполнение команды: %s\n", command)
	output, err := runCommand(client, command)
	if err != nil {
		log.Fatalf("Ошибка выполнения команды: %v", err)
	}

	fmt.Println("Результат выполнения команды:")
	fmt.Println(output)
}
