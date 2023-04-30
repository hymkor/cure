package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"math"
	"os"
	"regexp"
	"strings"
	"unicode/utf8"

	"golang.org/x/term"

	"github.com/mattn/go-colorable"

	"github.com/nyaosorg/go-windows-mbcs"
)

var ansiStrip = regexp.MustCompile("\x1B[^a-zA-Z]*[A-Za-z]")
var ansiOut = colorable.NewColorableStdout()

var bold = flag.Bool("b", false, "Use bold")
var screenWidth int
var screenHeight int

func splitLinesWithWidth(text string, screenWidth int, getWidth func(rune) int) (lines []string) {
	var buffer strings.Builder
	w := 0
	ansiStrip := false
	for _, c := range text {
		if ansiStrip {
			if ('A' <= c && c <= 'Z') || ('a' <= c && c <= 'z') {
				ansiStrip = false
			}
		} else if c == '\x1B' {
			ansiStrip = true
		} else {
			w1 := getWidth(c)
			if w+w1 >= screenWidth {
				lines = append(lines, buffer.String())
				buffer.Reset()
				w = 0
			}
			w += w1
		}
		buffer.WriteRune(c)
	}
	if buffer.Len() > 0 {
		lines = append(lines, buffer.String())
	}
	if len(lines) <= 0 {
		lines = []string{""}
	}
	return
}

func cat1(r io.Reader) error {
	sc := bufio.NewScanner(r)
	conin, err := newConin()
	if err != nil {
		return err
	}
	defer conin.Close()

	getWidth := newGetWidth(true)

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
		lines := splitLinesWithWidth(text, screenWidth, getWidth)
		for _, line := range lines {
			if count+1 >= screenHeight {
				io.WriteString(os.Stderr, "more>")
				ch, err := conin.getkey()
				if err != nil {
					return err
				}
				if ch == "\x03" {
					fmt.Fprintln(os.Stderr, "^C")
					return io.EOF
				}
				io.WriteString(os.Stderr, "\r     \b\b\b\b\b")
				if ch == "q" {
					return io.EOF
				} else if ch == "\r" {
					count--
				} else {
					count = 0
				}
				if *bold {
					io.WriteString(ansiOut, "\x1B[1m")
				}
			}
			fmt.Fprintln(ansiOut, line)
			count++
		}
	}
	return nil
}

func mains(args []string) error {
	count := 0

	var err error
	screenWidth, screenHeight, err = term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		return err
	}

	if !term.IsTerminal(int(os.Stdout.Fd())) {
		screenHeight = math.MaxInt32
	}

	for _, arg1 := range args {
		r, err := os.Open(arg1)
		if err != nil {
			return err
		}
		err = cat1(r)
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
		cat1(os.Stdin)
	}
	return nil
}

func main() {
	flag.Parse()
	if err := mains(flag.Args()); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
