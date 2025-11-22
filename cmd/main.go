package main

import (
	"log"

	"github.com/codingbot24-s/distributed-kv-store/internal/helper"
	"github.com/codingbot24-s/distributed-kv-store/internal/router"
)

func main() {
	err := helper.NewWal("Wal.log")
	if err != nil {
		log.Fatalf("error creating wal: %v", err)
	}
	helper.NewEngine()
	// start the read
	//err = helper.BuildState()
	//if err != nil {
	//	log.Fatalf("error building state: %v", err)
	//}
	router.StartRouter()

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
