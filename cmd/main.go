package main

import (
	storage "github.com/codingbot24-s/distributed-kv-store/internal"
)

func main() {
	w, err := storage.NewWal("Wal.log")
	if err != nil {
		panic(err)
	}
	//err = w.Append([]byte("hello wal"))
	//if err != nil {
	//	panic(err)
	//}

	w.Read()
	w.Close()
}
