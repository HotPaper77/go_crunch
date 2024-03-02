package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
)

var (
	inputDirectory string
)

func init() {
	flag.StringVar(&inputDirectory, "dir", ".", "Please type the directory to traverse")
	flag.Parse()
}

func main() {
	err := filepath.WalkDir(inputDirectory, processFiles)
	if err != nil {
		fmt.Printf("Could not traverse the directory %v, error: %v", inputDirectory, err)
	}
}

func processFiles(path string, d fs.DirEntry, err error) error {
	if err != nil {
		log.Fatal(err)
	}

	if d.IsDir() {
		return nil
	}

	rowsNum, err := countRows(filepath.Join(path))
	if err != nil {
		log.Fatalf("Could not process file:%v, error:%v", path, err)
		return err
	}

	fmt.Println(rowsNum)
	return nil
}

func countRows(path string) (int, error) {
	fd, err := os.Open(path)
	if err != nil {
		return 0, err
	}
	defer fd.Close()

	reader := csv.NewReader(fd)

	// Use Read instead of ReadAll to avoid unnecessary allocations
	var rowNumber int
	for {
		_, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			return 0, err
		}
		rowNumber++
	}

	return rowNumber, nil
}
