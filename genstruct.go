package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	fd "github.com/digisan/gotk/filedir"
	gio "github.com/digisan/gotk/io"
	"github.com/digisan/gotk/strs"
	jt "github.com/digisan/json-tool"
	lk "github.com/digisan/logkit"
)

func GenEntityPathVal(fpaths ...string) map[string]string {
	m := make(map[string]string)
	for _, fpath := range fpaths {
		if strs.HasAnySuffix(fpath, "class-link.json", "collection-entities.json") {
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
			path = strings.ReplaceAll(path, `.`, `[dot]`)
			val = strings.ReplaceAll(val.(string), `"`, `\"`)
			js += fmt.Sprintf(`"%s": "%s",`, path, val)
		}
		js = strings.TrimSuffix(js, ",") + "}"
		lk.FailOnErrWhen(!jt.IsValidStr(js), "%v @"+fpath, errors.New("invalid path-value json"))
		m[key.(string)] = js
	}
	return m
}

func DumpPathValue(idir, odname string) {
	osdir := filepath.Join(idir, odname)
	gio.MustCreateDir(osdir)
	fpaths, _, err := fd.WalkFileDir(idir, false)
	lk.FailOnErr("%v", err)
	for entity, js := range GenEntityPathVal(fpaths...) {
		lk.FailOnErr("%v", os.WriteFile(filepath.Join(osdir, entity+".json"), []byte(js), os.ModePerm))
	}
}
