package main

import (
	"errors"
	"flag"
	"os"
	"path/filepath"
	"strings"

	. "github.com/digisan/go-generics/v2"
	fd "github.com/digisan/gotk/filedir"
	gio "github.com/digisan/gotk/io"
	lk "github.com/digisan/logkit"
	"github.com/tidwall/gjson"
)

func init() {
	lk.Log2F(true, "./preproc.log")
	lk.WarnDetail(false)
}

//////////////////////////////////////////////////////////////////////////////////////

//
// 1) original => renamed --- cmd `go build -o rename`
//
func main() {

	var (
		dirOriEntPtr = flag.String("oed", "./data/original", "original entities json data directory")
		dirOriColPtr = flag.String("ocd", "./data/original/collections", "original collections json data directory")
		dirRnEntPtr  = flag.String("red", "./data/renamed", "renamed entities json data directory")
		dirRnColPtr  = flag.String("rcd", "./data/renamed/collections", "renamed collections json data directory")
	)

	flag.Parse()

	//////////////////////////////////////////////////////////////

	dirOriEnt, dirRnEnt := *dirOriEntPtr, *dirRnEntPtr

	gio.MustCreateDir(dirRnEnt)

	// clear destination dir for putting renamed file
	lk.FailOnErr("%v", fd.RmFilesIn(dirRnEnt, false, true, "json"))

	// make sure each file's name is its entity value
	FixFileName(dirOriEnt, dirRnEnt)

	//////////////////////////////////////////////////////////////

	dirOriCol, dirRnCol := *dirOriColPtr, *dirRnColPtr

	gio.MustCreateDir(dirRnCol)

	// clear destination dir for putting renamed file
	lk.FailOnErr("%v", fd.RmFilesIn(dirRnCol, false, true, "json"))

	// make sure each file's name is its entity value
	FixFileName(dirOriCol, dirRnCol)

	// <------------------------------------------------------------------------------------->

	mChk := map[string][]string{
		dirRnEnt: {"Element", "Object", "Abstract Element"},
		dirRnCol: {"Collection"},
	}

	for _, dir := range []string{dirRnEnt, dirRnCol} {
		fs, err := os.ReadDir(dir)
		lk.FailOnErr("%v", err)
		for _, f := range fs {
			if fname := f.Name(); strings.HasSuffix(fname, ".json") {
				fpath := filepath.Join(dir, fname)
				data, err := os.ReadFile(fpath)
				lk.FailOnErr("%v", err)
				lk.WarnOnErrWhen(NotIn(gjson.Get(string(data), "Metadata.Type").String(), mChk[dir]...), "%v@%s", errors.New("ERROR TYPE"), fpath)
			}
		}
	}

}

//////////////////////////////////////////////////////////////////////////////////////

//
// 2) renamed => out / err --- cmd: `go build -o preproc`
//
// func main() {

// 	dirInPtr := flag.String("in", "./data", "data directory")
// 	flag.Parse()

// 	dirIn := *dirInPtr

// 	out := filepath.Join(dirIn, "out")       // "out" is final output directory for ingestion
// 	errfolder := filepath.Join(dirIn, "err") // "err" is for incorrect format json dump into
// 	lk.FailOnErr("%v", os.RemoveAll(out))
// 	lk.FailOnErr("%v", os.RemoveAll(errfolder))

// 	Preproc(dirIn, out, errfolder)

// 	/////////////////////////////////////////////////////////////////////

// 	files, _, err := fd.WalkFileDir(out, false)
// 	lk.FailOnErr("%v", err)

// 	linkCol := LinkEntities(files...)

// 	js, err := Link2JSON(linkCol, "")
// 	lk.FailOnErr("%v", err)

// 	lk.FailOnErr("%v", os.WriteFile(filepath.Join(out, "class-link.json"), []byte(js), os.ModePerm))

// 	/////////////////////////////////////////////////////////////////////

// 	osdir := filepath.Join(out, "path_val")
// 	gio.MustCreateDir(osdir)
// 	fpaths, _, err := fd.WalkFileDir(out, false)
// 	lk.FailOnErr("%v", err)
// 	for entity, js := range GenEntityPathVal(fpaths...) {
// 		lk.FailOnErr("%v", os.WriteFile(filepath.Join(osdir, entity+".json"), []byte(js), os.ModePerm))
// 	}

// }
