package storage

import (
	"bufio"
	"fmt"
	"hash/crc32"
	"io"
	"os"
	"regexp"
	"strconv"
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
	sum := crc32.Checksum(data, crc32.MakeTable(crc32.Castagnoli))
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

func (w *Wal) Read() ([][]byte,error){
	if w.f == nil {
		return nil,fmt.Errorf("the file is not initialized")
	}
	reader := bufio.NewReader(w.f)

	_, err := w.f.Seek(0, 0)
	if err != nil {
		return nil,fmt.Errorf("error seeking file: %w", err)
	}

	// Compile regex once, not on every iteration
	re := regexp.MustCompile(`\[length: (\d+)\] \[checksum: (\d+)\] \[paylaod: (.*?)\]`)
	lineSlice := make([][]byte,0)
	for {
		bline, err := reader.ReadBytes('\n')
		if err != nil && err != io.EOF {
			return nil, fmt.Errorf("error reading file: %w", err)
		}
		if err == io.EOF {
			break
		}

		line := string(bline)
		matches := re.FindStringSubmatch(line)

		if len(matches) != 4 {
			return nil, fmt.Errorf("invalid log format: %s", line)
		}

		length, err := strconv.Atoi(matches[1])
		if err != nil {
			return nil, fmt.Errorf("invalid length value: %s", matches[1])
		}

		storedChecksum, err := strconv.ParseUint(matches[2], 10, 32)
		if err != nil {
			return nil, fmt.Errorf("invalid checksum value: %s", matches[2])
		}

		payload := matches[3]

		// compute new checksum
		computedChecksum := crc32.Checksum([]byte(payload), crc32.MakeTable(crc32.Castagnoli))

		// check
		if computedChecksum != uint32(storedChecksum) {
			return nil, fmt.Errorf("checksum mismatch for payload '%s': expected %d, got %d",
				payload, storedChecksum, computedChecksum)
		}

		if len(payload) != length {
			return nil, fmt.Errorf("length mismatch for payload '%s': expected %d, got %d",
				payload, length, len(payload))
		}

		fmt.Printf(" Verified: [length: %d] [checksum: %d] [payload: %s]\n", length, computedChecksum, payload)

		lineSlice = append(lineSlice,bline )
	}

	return lineSlice,nil
}
