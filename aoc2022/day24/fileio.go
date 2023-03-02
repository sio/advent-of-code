package main

import (
	"bufio"
	"io"
	"log"
	"os"
)

// Common logic for CharIterator and LineIterator
type fileIterator struct {
	file io.ReadCloser
	err  error
}

func (iter *fileIterator) openFile(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	iter.file = file
	return nil
}

func (iter *fileIterator) Close() error {
	if iter.file == nil {
		return nil
	}
	return iter.file.Close()
}

func (iter *fileIterator) Error() error {
	return iter.err
}

// Iterate over file by character
type CharIterator struct {
	fileIterator
	reader *bufio.Reader
	value  rune
}

func (iter *CharIterator) Open(filename string) error {
	err := iter.openFile(filename)
	if err != nil {
		return err
	}
	iter.reader = bufio.NewReader(iter.file)
	return nil
}

func (iter *CharIterator) Next() bool {
	if iter.file == nil || iter.reader == nil {
		return false
	}
	char, _, err := iter.reader.ReadRune()
	if err == io.EOF {
		return false
	}
	if err != nil {
		log.Printf("unexpected error while reading rune: %v", err)
		iter.err = err
		return false
	}
	iter.value = char
	return true
}

func (iter *CharIterator) Value() rune {
	return iter.value
}

// Iterate over file line by line
type LineIterator struct {
	fileIterator
	scanner *bufio.Scanner
}

func (iter *LineIterator) Open(filename string) error {
	err := iter.openFile(filename)
	if err != nil {
		return err
	}
	iter.scanner = bufio.NewScanner(iter.file)
	return nil
}

func (iter *LineIterator) Next() bool {
	if iter.file == nil || iter.scanner == nil {
		return false
	}
	if iter.scanner.Scan() {
		return true
	}
	iter.err = iter.scanner.Err()
	return false
}

func (iter *LineIterator) Value() string {
	return iter.scanner.Text()
}
