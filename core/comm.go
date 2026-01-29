package core

import "syscall"

type FDcomm struct {
	FD int
}

func (f FDcomm) Write(b []byte) (int, error) {
	return syscall.Write(f.FD, b)
}

func (f FDcomm) Read(b []byte) (int, error) {
	return syscall.Read(f.FD, b)
}
