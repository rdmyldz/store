package store

import (
	"bytes"
	"crypto/sha1"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

const rootDir = "data"

func Put(key []byte, r io.Reader) error {
	h, err := hashIt(key)
	if err != nil {
		return fmt.Errorf("in Put(): %w", err)
	}

	f, err := create(h)
	if err != nil {
		return fmt.Errorf("in Put(): %w", err)
	}

	_, err = io.Copy(f, r)
	if err != nil {
		return fmt.Errorf("in Put(): %w", err)
	}

	return nil
}

func PutBytes(key []byte, data []byte) error {
	r := bytes.NewReader(data)
	return Put(key, r)
}

func Get(key []byte) (io.Reader, error) {
	b, err := GetBytes(key)
	if err != nil {
		return nil, err
	}

	return bytes.NewReader(b), nil
}
func GetBytes(key []byte) ([]byte, error) {
	h, err := hashIt(key)
	if err != nil {
		return nil, fmt.Errorf("in GetBytes(): %w", err)
	}

	_, filePath := getPaths(h)
	f, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("in GetBytes(): %w", err)
	}
	defer f.Close()

	b, err := io.ReadAll(f)
	if err != nil {
		return nil, fmt.Errorf("in GetBytes(): %w", err)
	}

	return b, nil
}

func create(h []byte) (*os.File, error) {
	dirPath, filePath := getPaths(h)

	err := os.MkdirAll(dirPath, os.ModePerm)
	if err != nil {
		return nil, fmt.Errorf("in create(): %w", err)
	}

	return os.Create(filePath)
}

func hashIt(key []byte) ([]byte, error) {
	h := sha1.New()
	_, err := h.Write(key)
	if err != nil {
		return nil, fmt.Errorf("error in hashIt: %w", err)
	}

	return h.Sum(nil), nil
}

func getPaths(h []byte) (dirPath, filePath string) {
	dirname := fmt.Sprintf("%x", h[:2])
	filename := fmt.Sprintf("%x", h[2:])

	return filepath.Join(rootDir, dirname), filepath.Join(rootDir, dirname, filename)
}
