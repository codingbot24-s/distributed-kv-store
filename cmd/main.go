package main

import (
	"fmt"

	"github.com/codingbot24-s/distributed-kv-store/internal/helper"
)

func main() {
	w, err := helper.NewWal("Wal.log")
	if err != nil {
		panic(err)
	}
	_ = helper.NewEngine()
	// read all the entry from file
	entries, err := w.Read()
	if err != nil {
		panic(err)
	}
	//c := helper.Command{
	//	OP:    "set",
	//	Key:   "key",
	//	Value: "value",
	//}
	//err = helper.ApplyCommand(w, e, &c)
	//if err != nil {
	//	panic(err)
	//}
	// in the startup phase we need to read the wal file and get the
	// state back
	// startup phase
	for _, entry := range entries {
		// 1.entry will be in the byte so we need to decode it
		// 2. and get back the state for engine
		fmt.Println(string(entry))
	}

}
