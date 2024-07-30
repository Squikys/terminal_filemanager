package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
	"sort"

	"github.com/mattn/go-tty"
)

var nav int = 0
var start bool = false
var args [2]string

/*
	func IsHiddenFile(filename string) (bool, error) {
		pointer, err := syscall.UTF16PtrFromString(filename)
		if err != nil {
			return false, err
		}
		attributes, err := syscall.GetFileAttributes(pointer)
		if err != nil {
			return false, err
		}
		return attributes&syscall.FILE_ATTRIBUTE_HIDDEN != 0, nil
	}
*/
func opener(dir string) {
	var c *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		c = exec.Command("cmd", "/C", dir)

	default: //Mac & Linux
		c = exec.Command("rm", "-f", "/d/a.txt")
	}

	if err := c.Run(); err != nil {
		fmt.Println("Error: ", err)
	}
}
func navigation(dir []os.DirEntry) {
	tty, err := tty.Open()
	if err != nil {
		log.Fatal(err)
	}
	defer tty.Close()

	for {
		r, err := tty.ReadRune()
		if err != nil {
			log.Fatal(err)
		}
		if string(r) == "s" && nav < len(dir)-1 {
			nav++
		} else if string(r) == "w" && nav > 0 {
			nav--
		} else if string(r) == "o" {
			opener(string(args[1]) + dir[nav-1].Name())
		}

	}
}
func render(temp_nav int, dir []os.DirEntry) {
	if !start {
		for i := 0; i < len(dir); i++ {
			pointer := " "
			if nav == i {
				pointer = "*"
			}
			fmt.Printf("\033[%d;0H%s%s\n", i+1, pointer, dir[i].Name())
		}
		temp_nav = nav
		start = false
	}
	for {
		if nav != temp_nav {
			for i := 0; i < len(dir); i++ {
				pointer := " "
				if nav == i {
					pointer = "*"
				}
				fmt.Printf("\033[%d;0H%s%s\n", i+1, pointer, dir[i].Name())

			}
			temp_nav = nav
		}

	}
}
func main() {
	args[1] = os.Args[1]

	files, err := os.ReadDir(string(args[1]))
	sort.Slice(files, func(i, j int) bool {
		return true
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(len(files))
	fmt.Print("\033[2J")
	temp_nav := nav
	go navigation(files)
	render(temp_nav, files)

}
