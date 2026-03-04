package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"unicode/utf8"

	"github.com/nyaosorg/go-windows-mbcs"

	"github.com/tidwall/transform"
)

func nkf(src io.Reader) io.Reader {
	br := bufio.NewReader(src)
	return transform.NewTransformer(func() ([]byte, error) {
		line, err := br.ReadBytes('\n')
		if err != nil {
			return nil, err
		}
		if utf8.Valid(line) {
			return line, nil
		} else {
			text, err := mbcs.AtoU(line, mbcs.ACP)
			if err != nil {
				return nil, err
			}
			return []byte(text), nil
		}
	})
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
