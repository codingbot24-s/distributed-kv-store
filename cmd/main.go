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
	w,err := helper.GetWal()
	if err != nil {
		log.Fatalf("error getting wal: %v", err)
	}		

	e,err := helper.GetEngine()
	if err != nil {
		log.Fatalf("error getting engine : %v", err)
	}
	e.Replay(w)	
	e.Check()	
	router.StartRouter()	
}
