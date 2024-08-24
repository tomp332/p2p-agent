package storage

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
)

type File struct {
	FilePath   string
	buffer     *bytes.Buffer
	OutputFile *os.File
}

func (f *File) setFile() error {
	dir, singeFile := filepath.Split(filepath.Clean(f.FilePath))
	if singeFile == "" {
		return fmt.Errorf("path does not contain a file component")
	}
	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return err
	}
	file, err := os.Create(f.FilePath)
	if err != nil {
		return err
	}
	f.OutputFile = file
	return nil
}

func (f *File) Write(chunk []byte) error {
	if f.OutputFile == nil {
		err := f.setFile()
		if err != nil {
			return err
		}
	}
	_, err := f.OutputFile.Write(chunk)
	return err
}

func (f *File) Close() error {
	return f.OutputFile.Close()
}

func (f *File) Delete() error {
	if f.OutputFile == nil {
		return nil
	}
	err := os.Remove(f.FilePath)
	if err != nil {
		return err
	}
	return nil
}
