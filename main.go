package main

import (
	"flag"
	"os"
	"path/filepath"

	fd "github.com/digisan/gotk/filedir"
	lk "github.com/digisan/logkit"
)

func main() {

	lk.Log2F(true, "./preproc.log")

	dirOriPtr := flag.String("od", "./data/original", "original json data directory")
	dirRnPtr := flag.String("rd", "./data", "renamed json data directory")
	flag.Parse()

	dirOri, dirRn := *dirOriPtr, *dirRnPtr

	// make sure each file's name is its entity value
	FixFilename(dirOri, dirRn)

	// "out" is final output directory for ingestion; "err" is for incorrect format json dump into
	out := filepath.Join(dirRn, "out")
	Preproc(dirRn, out, filepath.Join(dirRn, "err"))

	// delete renamed file
	lk.FailOnErr("%v", fd.RmFilesIn(dirRn, false, true, "json"))

	/////////////////////////////////////////////////////////////////////

	files, _, err := fd.WalkFileDir(out, false)
	lk.FailOnErr("%v", err)

	linkCol := LinkEntities(files...)

	js, err := Link2JSON(linkCol, "")
	lk.FailOnErr("%v", err)

	lk.FailOnErr("%v", os.WriteFile(filepath.Join(out, "class-link.json"), []byte(js), os.ModePerm))
}
