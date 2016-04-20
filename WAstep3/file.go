package webAssign1

import (
	"os"
)

// read a file from the list.
func (gcs *gcsPhotos) read(filename string) (buf []byte, err error) {
	src, err := os.Open(filename)
	if err != nil {
		return
	}
	defer src.Close()
	size, err := src.Stat()
	if err != nil {
		return
	}
	buf = make([]byte, size.Size())
	_, err = src.Read(buf)
	if err != nil {
		return
	}
	return
}
