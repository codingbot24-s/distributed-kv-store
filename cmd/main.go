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
	e := helper.NewEngine()
	// read all the entry from file
	entries, err := w.Read()
	if err != nil {
		panic(err)
	}
	// startup phase we need to read the wal file and make the
	// engine state from it
	// TODO: we are correctly getting the all the entries just need to make the entry include
	// state from it with commands and key value
	for _, entry := range entries {
		// 1.entry will be in the byte so we need to decode it
		// 2. and get back the state for engine
		//TODO: is there a better way to get the command struct
		d := entry[46 : len(entry)-3]
		fmt.Println(string(d))
		c, err := helper.DecodeCommand(d)
		if err != nil {
			panic(err)
		}
		if err := e.Apply(c); err != nil {
			panic(err)
		}
		e.Check()
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

}
