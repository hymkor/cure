package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

func lnum(src io.Reader) io.Reader {
	in, out := io.Pipe()
	lnum := 0
	go func() {
		br := bufio.NewReader(src)
		for {
			line, err := br.ReadBytes('\n')
			if err != nil {
				out.CloseWithError(err)
				return
			}
			lnum++
			fmt.Fprintf(out, "%d: %s", lnum, line)
		}
	}()
	return in
}

func main() {
	sc := bufio.NewScanner(lnum(os.Stdin))
	for sc.Scan() {
		fmt.Println(sc.Text())
	}
	if err := sc.Err(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
