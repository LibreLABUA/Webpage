package main

import (
	"log"
	"os"
)

// this function allows to handle log
// without waiting the writing on disk
func InitLogger() chan<- string {
	c := make(chan string, 10)

	file, err := os.OpenFile(logfile, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		close(c)
		panic(err)
	}

	logger := log.New(file, "", log.Ldate)

	go func() {
		defer file.Close()
		for s := range c {
			logger.Println(s)
		}
	}()

	return c
}
