package main

import (
    "fmt"
    "flag"
    "io/ioutil"
    "path"
    "sync"
    "sort"
    "strings"
)

// Pair defines a key-value pair
type Pair struct {
  Key string
  Value int
}

// PairList defines a list of key-value pairs
type PairList []Pair

func (p PairList) Len() int { return len(p) }
func (p PairList) Less(i, j int) bool { return p[i].Value < p[j].Value }
func (p PairList) Swap(i, j int){ p[i], p[j] = p[j], p[i] }

// RankByWordCount is sorting given map by value
func RankByWordCount(wordFrequencies map[string]int) PairList{
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
    if (flag.NFlag() == 0) {
        flag.Usage()
        return
    }
    
    // Synchronization
    var waitgroup sync.WaitGroup
    waitgroup.Add(1)
    
    // Resulting map
    fileTypeMap := make(map[string]int)
    
    // Starting from the given directory
    go scan(*directory, &waitgroup, &fileTypeMap)
    
    // Waiting for completion
    waitgroup.Wait()
    fmt.Println();
    
    // Printing the resulting map
    var totalFilesCount int
    sortedList := RankByWordCount(fileTypeMap)
    for _, x := range sortedList {
        key := x.Key
        if (key == "") {
            key = "(no ext)"
        }
        fmt.Printf("%s\t\t=> %d\n", key, x.Value)
        totalFilesCount += x.Value
    }
    
    fmt.Println();
    fmt.Printf("Total files: %d, unique file types: %d", totalFilesCount, sortedList.Len())
}

func scan(directory string, waitgroup *sync.WaitGroup, fileTypeMap *map[string]int) {
    defer waitgroup.Done()
    fmt.Printf("scanning: %s\n", directory)
    
    // Scanning the directory
    files, _ := ioutil.ReadDir(directory)
    for _, file := range files {
            if (file.IsDir()) {
                waitgroup.Add(1)
                go scan(path.Join(directory, file.Name()), waitgroup, fileTypeMap)
            } else {
                extension := strings.ToLower(path.Ext(file.Name()))
                (*fileTypeMap)[extension]++
            }
    }
}
