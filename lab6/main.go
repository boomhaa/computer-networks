package main

import (
	bytes "bytes"
	"fmt"
	"github.com/jlaffaye/ftp"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"
)

func main() {
	c, err := ftp.Dial("students.yss.su:21", ftp.DialWithTimeout(5*time.Second))
	if err != nil {
		log.Fatal(err)
	}

	err = c.Login("ftpiu8", "3Ru7yOTA")
	if err != nil {
		log.Fatal(err)
	}

	command := ""
	path := ""
	curDir := "./"
	for command != "END" {
		fmt.Scan(&command)
		if command == "END" {
			if err := c.Quit(); err != nil {
				log.Fatal(err)
			}
			break
		} else if command == "GET" {
			fmt.Println("type a path")
			fmt.Scan(&path)
			r, err := c.Retr(path)
			if err != nil {
				panic(err)
			}
			defer r.Close()

			buf, err := ioutil.ReadAll(r)
			err = os.WriteFile(path, buf, 0666)
			if err != nil {
				log.Println(err)
			}
			fmt.Println("Done")
		} else if command == "PUSH" {
			fmt.Println("type a path")
			fmt.Scan(&path)
			text, err := os.ReadFile(path)
			if err != nil {
				log.Println(err)
			}
			byte := bytes.NewBuffer(text)
			err = c.Stor(path, byte)
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println("Done")
		} else if command == "MKDIR" {
			fmt.Println("type a path")
			fmt.Scan(&path)
			err = c.MakeDir(path)
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println("Done")
		} else if command == "CD" {
			fmt.Println("type a path")
			fmt.Scan(&path)
			curDir += path + "/"
			err = c.ChangeDir(path)
			if path == ".." {
				new_path := strings.Split(curDir, "/")
				curDir = strings.Join(new_path[:len(new_path)-3], "/")
			}
			fmt.Println(curDir)
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println("Done")
		} else if command == "RMDIR" {
			fmt.Println("type a path")
			fmt.Scan(&path)
			err = c.RemoveDir(path)
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println("Done")
		} else if command == "RM" {
			fmt.Println("type a path")
			fmt.Scan(&path)
			err = c.Delete(path)
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println("Done")
		} else if command == "LS" {
			list, err := c.List("./")
			if err != nil {
				fmt.Println(err)
			}
			for _, v := range list {
				fmt.Println(v.Name)
			}
			fmt.Println("Done")
		} else if command == "RMREC" {
			fmt.Println("type a path")
			fmt.Scan(&path)
			err = c.RemoveDirRecur(path)
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println("Done")
		}
	}

}
