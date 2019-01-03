package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"
	"unicode/utf8"

	"github.com/mattn/go-colorable"
	"github.com/mattn/go-runewidth"
	"github.com/mattn/go-tty"
	"github.com/zetamatta/go-texts/mbcs"
)

var ansiStrip = regexp.MustCompile("\x1B[^a-zA-Z]*[A-Za-z]")
var ansiOut = colorable.NewColorableStdout()

var bold = false
var screenWidth int
var screenHeight int

func cat1(r io.Reader, tty1 *tty.TTY) error {
	sc := bufio.NewScanner(r)
	count := 0
	for sc.Scan() {
		line := sc.Bytes()
		var text string
		if utf8.Valid(line) {
			text = string(line)
		} else {
			var err error
			text, err = mbcs.AtoU(line, mbcs.ACP)
			if err != nil {
				text = "ERROR: " + err.Error()
			}
		}
		width := runewidth.StringWidth(ansiStrip.ReplaceAllString(text, ""))
		lines := (width + screenWidth) / screenWidth
		for count+lines >= screenHeight {
			fmt.Fprint(os.Stderr, "more>")
			ch, err := tty1.ReadRune()
			if err != nil {
				return err
			}
			fmt.Fprint(os.Stderr, "\r     \b\b\b\b\b")
			if ch == 'q' {
				return io.EOF
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
	return nil
}

func main1() error {
	count := 0
	tty1, err := tty.Open()
	if err != nil {
		return err
	}
	defer tty1.Close()

	screenWidth, screenHeight, err = tty1.Size()
	if err != nil {
		return err
	}
	for _, arg1 := range os.Args[1:] {
		if arg1 == "-b" {
			bold = true
			continue
		} else if arg1 == "-h" {
			fmt.Println("CURE.exe : Color-Unicoded moRE")
			return nil
		}
		r, err := os.Open(arg1)
		if err != nil {
			return err
		}
		err = cat1(r, tty1)
		r.Close()
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}
		count++
	}
	if count <= 0 {
		cat1(os.Stdin, tty1)
	}
	return nil
}

func main() {
	if err := main1(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
