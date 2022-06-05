package main

import (
	"errors"
	"fmt"
	"os"
	"strings"

	jt "github.com/digisan/json-tool"
	lk "github.com/digisan/logkit"
)

func GenEntityPathVal(fpaths ...string) map[string]string {
	m := make(map[string]string)
	for _, fpath := range fpaths {
		if strings.HasSuffix(fpath, "class-link.json") {
			continue
		}
		data, err := os.ReadFile(fpath)
		lk.FailOnErr("%v", err)
		mPathVal, err := jt.Flatten(data)
		lk.FailOnErr("%v", err)
		key, ok := mPathVal["Entity"]
		lk.FailOnErrWhen(!ok, "%v @ "+fpath, errors.New("Entity Missing"))
		// make json
		js := "{"
		for path, val := range mPathVal {
			val = strings.ReplaceAll(val.(string), `"`, `\"`)
			js += fmt.Sprintf(`"%s": "%s",`, path, val)
		}
		js = strings.TrimSuffix(js, ",") + "}"
		lk.FailOnErrWhen(!jt.IsValidStr(js), "%v @"+fpath, errors.New("invalid path-value json"))
		m[key.(string)] = js
	}
	return m
}
