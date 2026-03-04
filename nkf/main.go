package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"unicode/utf8"

	"github.com/nyaosorg/go-windows-mbcs"
)

func nkf(src io.Reader) io.Reader {
	in, out := io.Pipe()
	go func() {
		br := bufio.NewReader(src)
		for {
			line, err := br.ReadBytes('\n')
			if err != nil {
				out.CloseWithError(err)
				return
			}
			if utf8.Valid(line) {
				_, err = out.Write(line)
				if err != nil {
					out.CloseWithError(err)
					return
				}
			} else {
				text, err := mbcs.AtoU(line, mbcs.ACP)
				if err != nil {
					out.CloseWithError(err)
					return
				}
				_, err = io.WriteString(out, text)
				if err != nil {
					out.CloseWithError(err)
					return
				}
			}
		}
	}()
	return in
}

func main() {
	sc := bufio.NewScanner(nkf(os.Stdin))
	for sc.Scan() {
		fmt.Println(sc.Text())
	}
	if err := sc.Err(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
