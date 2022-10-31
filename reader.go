package main

import (
	"bufio"
	"encoding/json"
	"io"
	"os"
)

type InputReader struct {
	file      *os.File
	reader    *bufio.Reader
	buffer    []byte
	separator byte
}

func NewReader(file *os.File, separator byte) *InputReader {
	fileReader := bufio.NewReader(file)
	return &InputReader{
		file:      file,
		reader:    fileReader,
		buffer:    []byte{},
		separator: separator,
	}
}

func (r *InputReader) Read() (map[string]any, error) {
	result := map[string]any{}
	content, err := r.reader.ReadBytes(r.separator)
	if err == io.EOF {
		if len(content) > 0 {
			r.buffer = append(r.buffer, content...)
		}
		return nil, nil
	}
	if err != nil {
		return result, err
	}
	err = json.Unmarshal(append(r.buffer, content...), &result)
	r.buffer = []byte{}
	return result, err
}
