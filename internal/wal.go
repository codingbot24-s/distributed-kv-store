package storage

import (
	"fmt"
	"os"
)

type Wal struct {
	f *os.File
}

func NewWal(path string) (*Wal, error) {
	f, err := os.OpenFile(path, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0644)
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

func (w *Wal) Append(data []byte) error {
	_, err := w.f.Write(data)
	if err != nil {
		return fmt.Errorf("error writing to the file: %w", err)
	}
	err = w.f.Sync()
	if err != nil {
		return fmt.Errorf("error syncing the file: %w", err)
	}

	return nil
}
