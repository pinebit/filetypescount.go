package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"path"
	"runtime"
	"sort"
	"strings"
	"sync"
)

// Pair defines a key-value pair
type Pair struct {
	Key   string
	Value int
}

// PairList defines a list of key-value pairs
type PairList []Pair

func (p PairList) Len() int           { return len(p) }
func (p PairList) Less(i, j int) bool { return p[i].Value < p[j].Value }
func (p PairList) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

// RankByWordCount is sorting given map by value
func RankByWordCount(wordFrequencies map[string]int) PairList {
	pl := make(PairList, len(wordFrequencies))
	i := 0
	for k, v := range wordFrequencies {
		pl[i] = Pair{k, v}
		i++
	}
	sort.Sort(sort.Reverse(pl))
	return pl
}

func main() {
	fmt.Println("A simple file type counter tool")

	// Handling command line flags
	directory := flag.String("dir", ".", "a directory to count files")
	flag.Parse()
	if flag.NFlag() == 0 {
		flag.Usage()
		return
	}

	// Synchronization
	var waitgroup sync.WaitGroup
	waitgroup.Add(1)

	// Resulting map and channels
	fileTypeMap := make(map[string]int)
	extensionsChan := make(chan string, 100)
	limitChan := make(chan int, runtime.NumCPU()*2)

	// Starting from the given directory
	go count(fileTypeMap, extensionsChan)
	go scan(*directory, &waitgroup, extensionsChan, limitChan)

	// Waiting for completion
	waitgroup.Wait()
	close(extensionsChan)
	fmt.Println()

	// Printing the resulting map
	var totalFilesCount int
	sortedList := RankByWordCount(fileTypeMap)
	for _, x := range sortedList {
		key := x.Key
		if key == "" {
			key = "(no ext)"
		}
		fmt.Printf("%s\t\t=> %d\n", key, x.Value)
		totalFilesCount += x.Value
	}

	fmt.Println()
	fmt.Printf("Total files: %d, unique file types: %d", totalFilesCount, sortedList.Len())
}

func scan(directory string, waitgroup *sync.WaitGroup, extensionsChan chan string, limitChan chan int) {
	defer waitgroup.Done()
	fmt.Printf("scanning: %s\n", directory)

	// Scanning the directory
	limitChan <- 1
	files, _ := ioutil.ReadDir(directory)
	for _, file := range files {
		if file.IsDir() {
			waitgroup.Add(1)
			go scan(path.Join(directory, file.Name()), waitgroup, extensionsChan, limitChan)
		} else {
			extension := strings.ToLower(path.Ext(file.Name()))
			extensionsChan <- extension
		}
	}
	<-limitChan
}

func count(fileTypeMap map[string]int, extensionsChan chan string) {
	for {
		extension, more := <-extensionsChan
		if more {
			fileTypeMap[extension]++
		} else {
			return
		}
	}
}
