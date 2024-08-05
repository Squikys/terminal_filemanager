package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/mattn/go-tty"
)

var nav int = 0
var files []os.DirEntry

func reducer(dir string) string {
	arr := strings.Split(dir, "\\")
	if len(arr) <= 1 || arr[1] == "" {
		return dir
	} else {
		var temp string = ""
		if arr[len(arr)-1] == "" {
			for i := 0; i < len(arr)-2; i++ {
				temp = temp + arr[i] + "\\"
			}
		} else {
			for i := 0; i < len(arr)-1; i++ {
				temp = temp + arr[i] + "\\"
			}
		}
		dir = temp
	}
	return dir
}
func next(path string, nextpath string) string {
	if string(path[len(path)-1]) != "\\" {
		path = path + "\\" + nextpath
		return path
	} else {
		path = path + nextpath
		return path
	}
}

func navigation(path string) {
	tty, err := tty.Open()
	if err != nil {
		log.Fatal(err)
	}
	defer tty.Close()
	render(files, path)
	for {
		r, err := tty.ReadRune()
		if err != nil {
			log.Fatal(err)
		}
		if string(r) == "s" && nav < len(files)-1 {
			nav++
			render(files, path)
		} else if string(r) == "w" && nav > 0 {
			nav--
			render(files, path)
		} else if string(r) == "q" {
			fmt.Print("\033[2J")
			path = reducer(path)
			files, _ = os.ReadDir(path)
			nav = 0
			render(files, path)
		} else if string(r) == "e" {
			if files[nav].Type().IsDir() {
				fmt.Print("\033[2J")
				path = next(path, files[nav].Name())
				files, _ = os.ReadDir(path)
				nav = 0
				render(files, path)
			}

		} else if string(r) == "v" {
			nav = 0
			cmd := exec.Command("cmd", "/C", "code", ".")
			_, err := cmd.Output()
			if err != nil {
				panic(err)
			}
		}

	}
}
func render(dir []os.DirEntry, path string) {

	fmt.Print("\033[1;0H current directory is")
	fmt.Print(path)

	for i := 0; i < len(dir); i++ {
		pointer := " "
		if nav == i {
			pointer = "*"
		}

		fmt.Printf("\033[%d;0H%s%s\n", i+3, pointer, dir[i].Name())

	}

}
func main() {
	args := os.Args[1]
	var err error
	files, err = os.ReadDir(string(args))
	path, _ := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Print("\033[2J")

	navigation(path)

}
