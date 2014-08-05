package main

import (
    "bufio"
    "fmt"
    "github.com/mfonda/simhash"
    "log"
    "os"
)

func readFile(filename string) ([]byte, int64, error) {
    file, err := os.OpenFile(filename, os.O_RDONLY, 0)
    if err != nil {
        return nil, -1, err
    }
    defer file.Close()

    stats, err := file.Stat()
    if err != nil {
        return nil, -1, err
    }

    var size int64 = stats.Size()
    bytes := make([]byte, size)

    buffer := bufio.NewReaderSize(file, 128*1024)
    _,err = buffer.Read(bytes)
    if err != nil {
        return nil, -1, err
    }

    return bytes, size, err
}

func main() {
    filename := os.Args[1]
    bytes, size, err := readFile(filename)
    if err != nil {
        log.Fatal(err)
    }

    hash := simhash.Simhash(simhash.NewWordFeatureSet(bytes))
    fmt.Printf("Simhash of %s (%v): %x\n", filename, size, hash)
}
