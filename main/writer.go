package main

import (
	"fmt"
	"log"
	"os"
	"time"
)

//Write file
func Write() {

	// If the file doesn't exist, create it, or append to the file
	f, err := os.OpenFile("./test.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0755)
	defer close(f)
	if err != nil {
		log.Fatal(err)
	}

	timeout := time.After(5 * time.Second)
	for {
		time.Sleep(500 * time.Millisecond)
		select {
		case <-timeout:
			fmt.Println("done writing")
			return
		default:
			if _, err := f.Write([]byte("appended some data\n")); err != nil {
				log.Fatal(err)
			}
		}
	}
}

func close(f *os.File) {
	f.Close()
}
