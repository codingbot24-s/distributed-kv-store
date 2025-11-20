package main

import (
	"log"

	"github.com/codingbot24-s/distributed-kv-store/internal/helper"
)

func main() {
	w, err := helper.NewWal("Wal.log")
	if err != nil {
		log.Fatalf("error creating wal: %v", err)
	}
	e := helper.NewEngine()
	// start the read
	err = helper.BuildState(w, e)
	if err != nil {
		log.Fatalf("error building state: %v", err)
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
