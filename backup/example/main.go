package main

import "badger-experiments/backup"

func main() {
	err := backup.Run()
	if err != nil {
		panic(err)
	}
}
