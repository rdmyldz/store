package store

import (
	"bytes"
	"crypto/sha1"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type Bundle struct{ rootDir string }

func NewBundle(root string) (Bundle, error) {
	if err := os.MkdirAll(root, os.ModePerm); err != nil {
		return Bundle{}, fmt.Errorf("store: %w", err)
	}
	return Bundle{rootDir: root}, nil
}

func (b Bundle) Put(key []byte, r io.Reader) error {
	h, err := hashIt(key)
	if err != nil {
		return fmt.Errorf("in Put(): %w", err)
	}

	f, err := create(b.rootDir, h)
	if err != nil {
		return fmt.Errorf("in Put(): %w", err)
	}

	_, err = io.Copy(f, r)
	if err != nil {
		return fmt.Errorf("in Put(): %w", err)
	}

	return nil
}

func (b Bundle) PutBytes(key []byte, data []byte) error {
	r := bytes.NewReader(data)
	return b.Put(key, r)
}

func (b Bundle) Get(key []byte) (io.Reader, error) {
	p, err := b.GetBytes(key)
	if err != nil {
		return nil, err
	}

	return bytes.NewReader(p), nil
}
func (b Bundle) GetBytes(key []byte) ([]byte, error) {
	h, err := hashIt(key)
	if err != nil {
		return nil, fmt.Errorf("in GetBytes(): %w", err)
	}

	_, filePath := getPaths(b.rootDir, h)
	f, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("in GetBytes(): %w", err)
	}
	defer f.Close()

	p, err := io.ReadAll(f)
	if err != nil {
		return nil, fmt.Errorf("in GetBytes(): %w", err)
	}

	return p, nil
}

func create(rootDir string, h []byte) (*os.File, error) {
	dirPath, filePath := getPaths(rootDir, h)

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

func getPaths(rootDir string, h []byte) (dirPath, filePath string) {
	dirname := fmt.Sprintf("%x", h[:2])
	filename := fmt.Sprintf("%x", h[2:])

	return filepath.Join(rootDir, dirname), filepath.Join(rootDir, dirname, filename)
}
