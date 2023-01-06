package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"math"
	"os"
	"regexp"
	"unicode/utf8"

	"github.com/mattn/go-colorable"
	"github.com/mattn/go-isatty"
	"github.com/mattn/go-runewidth"
	"github.com/mattn/go-tty"

	"github.com/nyaosorg/go-windows-mbcs"
)

var ansiStrip = regexp.MustCompile("\x1B[^a-zA-Z]*[A-Za-z]")
var ansiOut = colorable.NewColorableStdout()

var bold = flag.Bool("b", false, "Use bold")
var screenWidth int
var screenHeight int

func getkey() (rune, error) {
	tty1, err := tty.Open()
	if err != nil {
		return 0, err
	}
	defer tty1.Close()
	for {
		ch, err := tty1.ReadRune()
		if err != nil {
			return 0, err
		}
		if ch != 0 {
			return ch, nil
		}
	}
}

func cat1(r io.Reader) error {
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
			ch, err := getkey()
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
		if *bold {
			fmt.Fprint(ansiOut, "\x1B[1m")
		}
		fmt.Fprintln(ansiOut, text)
		count += lines
	}
	return nil
}

func main1(args []string) error {
	count := 0
	tty1, err := tty.Open()
	if err != nil {
		return err
	}
	screenWidth, screenHeight, err = tty1.Size()
	tty1.Close()
	if err != nil {
		return err
	}

	if !isatty.IsTerminal(os.Stdout.Fd()) {
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
	if err := main1(flag.Args()); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
