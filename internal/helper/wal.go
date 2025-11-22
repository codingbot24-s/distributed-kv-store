package helper

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
	f     *os.File
	Index int64
	Term  int64
}

var defaultWal *Wal

func NewWal(path string) error {
	f, err := os.OpenFile(path, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return fmt.Errorf("error opening the file: %w", err)
	}
	// store the created WAL in the package-level variable so GetWal can return it
	defaultWal = &Wal{f: f}
	return nil
}

func GetWal() (*Wal, error) {
	if defaultWal == nil {
		return nil, fmt.Errorf("wal is not initialized")
	}
	return defaultWal, nil
}

func (w *Wal) Close() error {
	err := w.f.Close()
	if err != nil {
		return fmt.Errorf("error closing the file: %w", err)
	}

	return nil
}

// TODO: change to add a log entry
func (w *Wal) Append(data []byte) error {
	//TODO: error in checksum because sum is different with only payload we need tp pass index and current term in byte so checksum will pass
	sum := crc32.Checksum(data, crc32.MakeTable(crc32.Castagnoli))

	line := fmt.Sprintf("[length: %d] [checksum: %d] [paylaod: %s] \n", len(data), sum,
		string(data))
	_, err := w.f.Write([]byte(line))
	if err != nil {
		return fmt.Errorf("error writing to the file: %w", err)
	}
	err = w.f.Sync()
	if err != nil {
		return fmt.Errorf("error syncing the file: %w", err)
	}
	w.Index++
	return nil
}

func (w *Wal) getTerm() int64 {
	return w.Term
}

func (w *Wal) getIndex() int64 {
	return w.Index
}

// do we need to return whole entry can we just return the payload and make
// TODO: before read would return the whole entry but now it is returning just theand
// not len and sum
func (w *Wal) Read() ([][]byte, error) {
	if w.f == nil {
		return nil, fmt.Errorf("the file is not initialized")
	}
	reader := bufio.NewReader(w.f)

	_, err := w.f.Seek(0, 0)
	if err != nil {
		return nil, fmt.Errorf("error seeking file: %w", err)
	}

	// Compile regex once, not on every iteration
	re := regexp.MustCompile(`\[length: (\d+)\] \[checksum: (\d+)\] \[paylaod: (.*?)\]`)
	lineSlice := make([][]byte, 0)
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
		// TODO: read error is here we need to get only command not term and index
		// we can split by command and then try to compute ?
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

		lineSlice = append(lineSlice, bline)
	}

	return lineSlice, nil
}
