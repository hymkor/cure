package main

import (
	"fmt"
	"syscall"

	"golang.org/x/term"
)

type Conin struct {
	in _Handle
}

func newConin() (*Conin, error) {
	in, err := syscall.Open(_ConIn, syscall.O_RDWR, 0)

	if err != nil {
		return nil, fmt.Errorf("newConin: %w", err)
	}
	return &Conin{in: in}, nil
}

func (C *Conin) Close() error {
	return syscall.Close(C.in)
}

func (C *Conin) getkey() (string, error) {
	stdin := int(C.in)
	state, err := term.MakeRaw(stdin)
	if err != nil {
		return "", err
	}
	defer term.Restore(stdin, state)

	for {
		var buffer [256]byte
		n, err := syscall.Read(C.in, buffer[:])
		if err != nil {
			return "", fmt.Errorf("(Conin.in) Read: %w", err)
		}
		if n > 0 {
			return string(buffer[:n]), nil
		}
	}
}
