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

func write(wg *sync.WaitGroup) {
	go func() {
		defer wg.Done()
		//Write()
	}()
}

//Tail structure
type Tail struct {
	lk sync.Mutex
}

func read(wg *sync.WaitGroup, ch chan string) {
	go func() {
		defer wg.Done()
		//defer close(ch)
		fmt.Println("reading file")
		file, err := os.Open("/go/test.log")
		fmt.Println("file opened for reading")
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
		tail := &Tail{}

		reader := bufio.NewReader(file)
		//timeout := time.After(5 * time.Second)
		count := 1
		for {
			select {
			//case <-timeout:
			//fmt.Println("done reading")
			//return
			default:
				tail.lk.Lock()
				text, err := reader.ReadString('\n')
				tail.lk.Unlock()
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
		//timeout := time.After(5 * time.Second)
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
				//fmt.Println("[", line, "]")
				/*case <-timeout:
				print(&buffer)
				return*/
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
	wg.Add(3)
	write(&wg)
	read(&wg, logCh)
	process(&wg, logCh)
	wg.Wait()
}
