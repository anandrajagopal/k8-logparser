package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
	"sync"
	"time"

	yaml "gopkg.in/yaml.v2"
)

//Tail structure
type Tail struct {
	Fn      string `yaml:"fileName"`
	Pattern string `yaml:"newLinePattern"`
}

func read(applog *Tail, wg *sync.WaitGroup, printChan chan string) {
	processChan := make(chan string)
	process(applog, wg, processChan, printChan)
	go func() {
		defer wg.Done()
		fmt.Println("reading file", applog.Fn)
		file, err := os.Open(applog.Fn)
		if err != nil {
		L:
			for {
				select {
				//wait until the file is available
				case <-time.After(5 * time.Second):
					fmt.Println("looking for file", applog.Fn)
					file, err = os.Open(applog.Fn)
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
					processChan <- text
				}
			}

		}
	}()
}

func process(applog *Tail, wg *sync.WaitGroup, processChan <-chan string, printChan chan string) {
	var buffer bytes.Buffer
	pattern := applog.Pattern
	wg.Add(1)
	go func() {
		defer wg.Done()
		count := 1
		for {
			select {
			case line := <-processChan:
				matched, err := regexp.MatchString(pattern, line)
				if err != nil {
					panic(err)
				}
				if matched {
					//print(&buffer)
					//buffer.Reset()
					printChan <- buffer.String()
					buffer.Reset()
					count = 1
				}
				buffer.WriteString(strconv.Itoa(count) + "--> " + line)
				count++
			}
		}
	}()
}

func printFromChan(wg *sync.WaitGroup, printChan chan string) {
	go func() {
		defer wg.Done()
		for {
			select {
			case line := <-printChan:
				fmt.Print("{{", line, "}}")
			}
		}
	}()
}

func print(r io.Reader) {
	fmt.Println("{{", r, "}}")
	//io.Copy(os.Stdout, r)
}

func main() {
	bytes, err := ioutil.ReadFile("/etc/config/config")
	if err != nil {
		panic("Cannot read file")
	}
	logs := []*Tail{}
	err = yaml.Unmarshal(bytes, &logs)
	if err != nil {
		panic(err)
	}
	var wg sync.WaitGroup
	logCh := make(chan string)
	wg.Add(2)
	for _, log := range logs {
		read(log, &wg, logCh)
	}
	printFromChan(&wg, logCh)
	wg.Wait()
}
