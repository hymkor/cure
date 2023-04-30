package main

import (
	"syscall"
)

type _Handle = syscall.Handle

const (
	_ConIn = "CONIN$"
)
