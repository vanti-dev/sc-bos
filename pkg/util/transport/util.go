package transport

import (
	"io"
)

func readWriteCloserFuncs(read, write readWrite, close func() error) io.ReadWriteCloser {
	return &readWriteCloser{
		read:  read,
		write: write,
		close: close,
	}
}

type readWriteCloser struct {
	read  readWrite
	write readWrite
	close func() error
}

func (r *readWriteCloser) Read(p []byte) (n int, err error) {
	return r.read(p)
}

func (r *readWriteCloser) Write(p []byte) (n int, err error) {
	return r.write(p)
}

func (r *readWriteCloser) Close() error {
	return r.close()
}

type readWrite func(p []byte) (n int, err error)
