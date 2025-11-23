package helper

import (
	"bufio"
	"encoding/json"
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
// cheksum whole data before we are computing the checksum of only command now whole payloa so read and write dosnt generate the diffrent checksum
func (w *Wal) Append(data []byte) error {
	// Remove the slicing - checksum the entire JSON payload
	sum := crc32.Checksum(data, crc32.MakeTable(crc32.Castagnoli))

	// Fix the typo: paylaod -> payload
	line := fmt.Sprintf("[length: %d] [checksum: %d] [payload: %s]\n",
		len(data), sum, string(data))

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

func (w *Wal) GetTerm() int64 {
	return w.Term
}

func (w *Wal) GetIndex() int64 {
	return w.Index
}


func (w *Wal) Read() ([]*LogEntry, error) {
	if w.f == nil {
		return nil, fmt.Errorf("the file is not initialized")
	}

	reader := bufio.NewReader(w.f)
	_, err := w.f.Seek(0, 0)
	if err != nil {
		return nil, fmt.Errorf("error seeking file: %w", err)
	}

	// Match both "payload" and "paylaod" for compatibility
	re := regexp.MustCompile(`\[length: (\d+)\] \[checksum: (\d+)\] \[pay(?:loa|la)d: (.*?)\]`)

	var entries []*LogEntry 

	for {
		line, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("error reading file: %w", err)
		}

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

		payload := []byte(matches[3])

		// Verify checksum on the JSON payload
		computedChecksum := crc32.Checksum(payload, crc32.MakeTable(crc32.Castagnoli))
		if computedChecksum != uint32(storedChecksum) {
			return nil, fmt.Errorf("checksum mismatch: expected %d, got %d",
				storedChecksum, computedChecksum)
		}

		if len(payload) != length {
			return nil, fmt.Errorf("length mismatch: expected %d, got %d",
				length, len(payload))
		}

		// Parse JSON into LogEntry struct
		var logEntry LogEntry
		if err := json.Unmarshal(payload, &logEntry); err != nil {
			return nil, fmt.Errorf("error unmarshaling entry: %w", err)
		}

		entries = append(entries, &logEntry)
	}

	return entries, nil
}
