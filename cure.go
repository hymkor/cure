package main

import "bufio"
import "io"
import "os"
import "regexp"
import "fmt"

import "github.com/zetamatta/nyagos/conio"
import "github.com/shiena/ansicolor"

var ansiStrip = regexp.MustCompile("\x1B[^a-zA-Z]*[A-Za-z]")
var ansiOut = ansicolor.NewAnsiColorWriter(os.Stdout)

var bold = false
var screenWidth int
var screenHeight int

func cat1(r io.Reader) bool {
	scanner := bufio.NewScanner(r)
	count := 0
	for scanner.Scan() {
		width := conio.GetStringWidth(ansiStrip.ReplaceAllString(scanner.Text(), ""))
		lines := (width + screenWidth) / screenWidth
		for count+lines > screenHeight {
			fmt.Print("more>")
			ch := conio.GetCh()
			fmt.Print("\r     \b\b\b\b\b")
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
		fmt.Fprintln(ansiOut, scanner.Text())
		count += lines
	}
	return true
}

func main() {
	count := 0
	screenWidth, screenHeight = conio.GetScreenSize()
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
