package main

import (
	"os"
	"path/filepath"
	"testing"

	fd "github.com/digisan/gotk/filedir"
	gio "github.com/digisan/gotk/io"
	lk "github.com/digisan/logkit"
)

func TestGenEntityPathVal(t *testing.T) {
	out := "./data/out"
	osdir := filepath.Join(out, "path_val")
	gio.MustCreateDir(osdir)
	fpaths, _, err := fd.WalkFileDir(out, false)
	lk.FailOnErr("%v", err)
	for entity, js := range GenEntityPathVal(fpaths...) {
		lk.FailOnErr("%v", os.WriteFile(filepath.Join(osdir, entity+".json"), []byte(js), os.ModePerm))
	}
}
