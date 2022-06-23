package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	lk "github.com/digisan/logkit"
	"github.com/tidwall/gjson"
)

func DumpCollection(dir, ofname string) {
	mColEntities := make(map[string][]string)
	fis, err := os.ReadDir(dir)
	lk.FailOnErr("%v", err)
	for _, fi := range fis {
		if fname := fi.Name(); strings.HasSuffix(fname, ".json") {
			bytesJS, err := os.ReadFile(filepath.Join(dir, fname))
			lk.FailOnErr("%v", err)
			js := string(bytesJS)
			if r := gjson.Get(js, "Collections"); r.IsArray() {
				for i := 0; i < len(r.Array()); i++ {
					colName := gjson.Get(js, fmt.Sprintf("Collections.%d.Name", i)).String()
					mColEntities[colName] = append(mColEntities[colName], strings.TrimSuffix(fname, ".json"))
				}
			}
		}
	}
	bytesJS, err := json.Marshal(mColEntities)
	lk.FailOnErr("%v", err)
	lk.FailOnErr("%v", os.WriteFile(filepath.Join(dir, ofname), bytesJS, os.ModePerm))
}
