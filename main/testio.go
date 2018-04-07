package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"regexp"
	"sync"
	"time"
)

//Tail structure
type Tail struct {
	lk      sync.Mutex
	fn      string
	pattern string
}

func read(wg *sync.WaitGroup, ch chan string) {
	go func() {
		defer wg.Done()
		//defer close(ch)
		fmt.Println("reading file")
		file, err := os.Open("/go/test.log")
		if err != nil {
		L:
			for {
				select {
				//wait until the file is available
				case <-time.After(5 * time.Second):
					fmt.Println("looking for file", "/go/test.log")
					file, err = os.Open("/go/test.log")
					if err == nil {
						break L
					}
				}
			}
		}
		fmt.Println("file opened for reading")

		//tail := &Tail{}

		reader := bufio.NewReader(file)
		count := 1
		for {
			select {
			default:
				//tail.lk.Lock()
				text, err := reader.ReadString('\n')
				//tail.lk.Unlock()
				if err == io.EOF {
					//fmt.Println("EOF reached")
					//break
				}
				if err != nil {
					//panic("unable to read file")
				} else {
					ch <- string(count) + text
					count++
				}
			}

		}
	}()
}

func process(wg *sync.WaitGroup, logCh <-chan string) {
	var buffer bytes.Buffer
	pattern := `\d{4}-\d{2}-\d{2}\s\d{2}`
	go func() {
		defer wg.Done()
		for {
			select {
			case line := <-logCh:
				matched, err := regexp.MatchString(pattern, line)
				if err != nil {
					panic(err)
				}
				if matched {
					print(&buffer)
					buffer.Reset()
				}
				buffer.WriteString(line)
			}
		}
	}()
}

func print(r io.Reader) {
	fmt.Println("{{", r, "}}")
	//io.Copy(os.Stdout, r)
}

func main() {
	var wg sync.WaitGroup
	logCh := make(chan string)
	wg.Add(2)
	read(&wg, logCh)
	process(&wg, logCh)
	wg.Wait()
}
