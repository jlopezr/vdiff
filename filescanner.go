package main

import (
	"bufio"
	"strings"
)

type FileScanner struct {
	scanner *bufio.Scanner
	fileName string
	hash string
	hasMore bool
}

func NewFileScanner(scanner *bufio.Scanner) *FileScanner {
	m := new(FileScanner)
	m.scanner = scanner
	m.hasMore = true
	return m
}

func (f *FileScanner) Scan() bool {

	for {
		result := f.scanner.Scan()
		if !result {
			f.hasMore = false
			return false
		}

		line := f.scanner.Text()

		if strings.HasSuffix(line, "is a dir") {
			continue
		}

		f.hash = line[len(line)-32:]
		f.fileName = line[0:len(line)-33]
		f.hasMore = true
		return true
	}
}

func (f *FileScanner) HasMore() bool {
	return f.hasMore
}

func (f *FileScanner) Text() string {
	return f.scanner.Text()
}

func (f *FileScanner) Err() error {
	return f.scanner.Err()
}

func (f *FileScanner) Hash() string {
	return f.hash
}

func (f *FileScanner) FileName() string {
	return f.fileName
}