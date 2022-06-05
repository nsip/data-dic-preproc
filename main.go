package main

import (
	"flag"
	"os"
	"path/filepath"

	fd "github.com/digisan/gotk/filedir"
	gio "github.com/digisan/gotk/io"
	lk "github.com/digisan/logkit"
)

func init() {
	lk.Log2F(true, "./preproc.log")
}

//////////////////////////////////////////////////////////////////////////////////////

//
// 1) original => renamed --- cmd `go build -o rename`
//
// func main() {

// 	dirOriPtr := flag.String("od", "./data/original", "original json data directory")
// 	dirRnPtr := flag.String("rd", "./data", "renamed json data directory")
// 	flag.Parse()

// 	dirOri, dirRn := *dirOriPtr, *dirRnPtr

// 	// clear destination dir for putting renamed file
// 	lk.FailOnErr("%v", fd.RmFilesIn(dirRn, false, true, "json"))

// 	// make sure each file's name is its entity value
// 	FixFileName(dirOri, dirRn)
// }

//////////////////////////////////////////////////////////////////////////////////////

//
// 2) renamed => out / err --- cmd: `go build -o preproc`
//
func main() {

	dirInPtr := flag.String("in", "./data", "data directory")
	flag.Parse()

	dirIn := *dirInPtr

	out := filepath.Join(dirIn, "out")       // "out" is final output directory for ingestion
	errfolder := filepath.Join(dirIn, "err") // "err" is for incorrect format json dump into
	lk.FailOnErr("%v", os.RemoveAll(out))
	lk.FailOnErr("%v", os.RemoveAll(errfolder))

	Preproc(dirIn, out, errfolder)

	/////////////////////////////////////////////////////////////////////

	files, _, err := fd.WalkFileDir(out, false)
	lk.FailOnErr("%v", err)

	linkCol := LinkEntities(files...)

	js, err := Link2JSON(linkCol, "")
	lk.FailOnErr("%v", err)

	lk.FailOnErr("%v", os.WriteFile(filepath.Join(out, "class-link.json"), []byte(js), os.ModePerm))

	/////////////////////////////////////////////////////////////////////

	osdir := filepath.Join(out, "path_val")
	gio.MustCreateDir(osdir)
	fpaths, _, err := fd.WalkFileDir(out, false)
	lk.FailOnErr("%v", err)
	for entity, js := range GenEntityPathVal(fpaths...) {
		lk.FailOnErr("%v", os.WriteFile(filepath.Join(osdir, entity+".json"), []byte(js), os.ModePerm))
	}

}
