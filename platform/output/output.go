package output

import (
	"log"
	"os"
)

// TODO: EventWrite struct instead of Writer to use custom formatting of output information

type Writer struct {
	file *os.File
}

func (ew *Writer) Write(data []byte) (int, error) {
	return ew.file.Write(data)
}

func (ew *Writer) Close() {
	ew.file.Close()
}

func MustMakeWriter(filename string) *Writer {
	file, err := os.Create(filename)
	if err != nil {
		log.Fatalf("failed to create output file: %v", err)
	}
	return &Writer{file: file}
}
