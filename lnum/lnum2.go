package main

import (
	"bufio"
	"fmt"
	"io"
	"os"

	"github.com/tidwall/transform"
)

func lnum(src io.Reader) io.Reader {
	br := bufio.NewReader(src)
	lnum := 0
	return transform.NewTransformer(func() ([]byte, error) {
		line, err := br.ReadString('\n')
		if err != nil {
			return nil, err
		}
		lnum++
		return []byte(fmt.Sprintf("%d: %s", lnum, line)), nil
	})
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
