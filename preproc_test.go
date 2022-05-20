package main

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestFixFilename(t *testing.T) {
	FixFilename("./data/original", "./data")
}

func TestPreproc(t *testing.T) {

	filepath.Walk("./", func(path string, info fs.FileInfo, err error) error {
		var e error
		if strings.HasSuffix(path, ".json") {
			e = os.Remove(path)
		}
		return e
	})

	Preproc("./data", "./data/out", "./data/err")
}

func TestRmP(t *testing.T) {
	data, err := os.ReadFile("./data/Education support.json")
	if err != nil {
		panic(err)
	}
	removedP := rmPtag(string(data))
	err = os.WriteFile("./data/out.json", []byte(removedP), os.ModePerm)
	if err != nil {
		panic(err)
	}
}
