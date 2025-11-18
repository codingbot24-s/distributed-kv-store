package storage

import (
	"fmt"
	"os"
)

type Wal struct {
	f *os.File
}

func NewWal(path string) (*Wal, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("error opening the file: %w", err)
	}

	return &Wal{f: f}, nil
}
func (w *Wal) Close() {
	err := w.f.Close()
	if err != nil {
		fmt.Printf("error closing the file: %v", err)
	}
}
