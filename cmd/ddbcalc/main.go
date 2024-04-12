package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/ryeguard/ddbcalc"
)

var fileFlag = flag.String("file", "", "The file path and name to read. Must be a json file. Required.")

func readJSON(file string) (map[string]interface{}, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, fmt.Errorf("open: %w", err)
	}
	defer f.Close()

	b, err := io.ReadAll(f)
	if err != nil {
		return nil, fmt.Errorf("read: %w", err)
	}

	var data map[string]interface{}

	err = json.Unmarshal(b, &data)
	if err != nil {
		return nil, fmt.Errorf("unmarshal: %w", err)
	}

	return data, nil
}

func main() {
	flag.Parse()

	if *fileFlag == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	data, err := readJSON(*fileFlag)
	if err != nil {
		panic(fmt.Sprintf("readJSON: %v", err))
	}

	size, err := ddbcalc.StructSizeInBytes(data)
	if err != nil {
		panic(fmt.Sprintf("StructSizeInBytes: %v", err))
	}

	fmt.Printf("The resulting DynamoDB item size is %d bytes\n", size)
}
