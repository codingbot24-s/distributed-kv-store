package storage

import (
	"bufio"
	"bytes"
	"fmt"
	"hash/crc32"
	"io"
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
func (w *Wal) Close() error {
	err := w.f.Close()
	if err != nil {
		return fmt.Errorf("error closing the file: %w", err)
	}

	return nil
}

func (w *Wal) Append(data []byte) error {
	sum := crc32.Checksum(data, crc32.IEEETable)
	line := fmt.Sprintf("[length: %d] [checksum: %d] [paylaod: %s] \n", len(data), sum, string(data))
	_, err := w.f.Write([]byte(line))
	if err != nil {
		return fmt.Errorf("error writing to the file: %w", err)
	}
	err = w.f.Sync()
	if err != nil {
		return fmt.Errorf("error syncing the file: %w", err)
	}

	return nil
}

type Oneline struct {
	Length  int32
	sum     int32
	Payload string
}

func (w *Wal) Read() error {
	if _, err := w.f.Seek(0, io.SeekStart); err != nil {
		return fmt.Errorf("error seeking the file: %w", err)
	}
	scanner := bufio.NewScanner(w.f)
	for scanner.Scan() {
		line := scanner.Bytes()
		_ = bytes.TrimSpace(line)
		//TODO: how can we get the checksum

	}

	return nil
}
