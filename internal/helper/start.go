package helper

import "fmt"

func BuildState() error {
	w, err := GetWal()
	if err != nil {
		return fmt.Errorf("error getting wal: %v", err)
	}
	e, err := GetEngine()
	if err != nil {
		return fmt.Errorf("error getting engine: %v", err)
	}
	entries, err := w.Read()
	if err != nil {
		return fmt.Errorf("error in reading wal %w", err)
	}
	for _, entry := range entries {
		
		//TODO: is there a better way to get the command struct
		d := entry[77 : len(entry)-5]
		fmt.Println(string(d))
		c, err := DecodeCommand(d)
		if err != nil {
			return fmt.Errorf("error in decoding command %w", err)
		}
		if err := e.Apply(c); err != nil {
			return fmt.Errorf("error in applying command %w", err)
		}
		e.Check()
	}

	return nil
}
