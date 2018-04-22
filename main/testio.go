package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
	"sync"
	"time"
)

//Tail structure
type Tail struct {
	lk      sync.Mutex
	fn      string
	pattern string
}

func read(applog *Tail, wg *sync.WaitGroup, ch chan string) {
	go func() {
		defer wg.Done()
		//defer close(ch)
		fmt.Println("reading file", applog.fn)
		file, err := os.Open(applog.fn)
		if err != nil {
		L:
			for {
				select {
				//wait until the file is available
				case <-time.After(5 * time.Second):
					fmt.Println("looking for file", applog.fn)
					file, err = os.Open(applog.fn)
					if err == nil {
						break L
					}
				}
			}
		}
		fmt.Println("file opened for reading")

		//tail := &Tail{}

		reader := bufio.NewReader(file)
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
					ch <- text
				}
			}

		}
	}()
}

func process(applog *Tail, wg *sync.WaitGroup, logCh <-chan string) {
	var buffer bytes.Buffer
	pattern := applog.pattern
	go func() {
		defer wg.Done()
		count := 1
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
					count = 1
				}
				buffer.WriteString(strconv.Itoa(count) + "--> " + line)
				count++
			}
		}
	}()
}

func print(r io.Reader) {
	fmt.Println("{{", r, "}}")
	//io.Copy(os.Stdout, r)
}

func main() {
	applog := &Tail{
		fn:      `/etc/log/output.log`,
		pattern: `\d{4}-\d{2}-\d{2}\s\d{2}`,
	}
	var wg sync.WaitGroup
	logCh := make(chan string)
	wg.Add(2)
	read(applog, &wg, logCh)
	process(applog, &wg, logCh)
	wg.Wait()
}
