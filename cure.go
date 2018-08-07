package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"

	"github.com/mattn/go-colorable"
	"github.com/mattn/go-runewidth"
	"github.com/zetamatta/go-console/screenbuffer"
	"github.com/zetamatta/go-getch"
	"github.com/zetamatta/go-texts/mbcs"
)

var ansiStrip = regexp.MustCompile("\x1B[^a-zA-Z]*[A-Za-z]")
var ansiOut = colorable.NewColorableStdout()

var bold = false
var screenWidth int
var screenHeight int

func cat1(r io.Reader) bool {
	scanner := bufio.NewScanner(mbcs.NewAutoDetectReader(r, mbcs.ConsoleCP()))
	count := 0
	for scanner.Scan() {
		text := scanner.Text()
		width := runewidth.StringWidth(ansiStrip.ReplaceAllString(text, ""))
		lines := (width + screenWidth) / screenWidth
		for count+lines >= screenHeight {
			fmt.Fprint(os.Stderr, "more>")
			ch := getch.Rune()
			fmt.Fprint(os.Stderr, "\r     \b\b\b\b\b")
			if ch == 'q' {
				return false
			} else if ch == '\r' {
				count--
			} else {
				count = 0
			}
		}
		if bold {
			fmt.Fprint(ansiOut, "\x1B[1m")
		}
		fmt.Fprintln(ansiOut, text)
		count += lines
	}
	return true
}

func main() {
	count := 0
	screenWidth, screenHeight = csbi.GetConsoleScreenBufferInfo().ViewSize()
	for _, arg1 := range os.Args[1:] {
		if arg1 == "-b" {
			bold = true
			continue
		} else if arg1 == "-h" {
			fmt.Println("CURE.exe : Color-Unicoded moRE")
			return
		}
		r, err := os.Open(arg1)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			return
		}
		if !cat1(r) {
			return
		}
		r.Close()
		count++
	}
	if count <= 0 {
		cat1(os.Stdin)
	}
}
