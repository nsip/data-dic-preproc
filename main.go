package main

import (
	"flag"

	// . "github.com/digisan/go-generics/v2"
	fd "github.com/digisan/gotk/filedir"
	gio "github.com/digisan/gotk/io"
	lk "github.com/digisan/logkit"
)

func init() {
	lk.Log2F(true, "./preproc.log")
	lk.WarnDetail(false)
}

//////////////////////////////////////////////////////////////////////////////////////

//
// 1) original => renamed --- cmd `go build -o rename`
//
// func main() {

// 	var (
// 		dirOriEntPtr = flag.String("oed", "./data/original", "original entities json data directory")
// 		dirOriColPtr = flag.String("ocd", "./data/original/collections", "original collections json data directory")
// 		dirRnEntPtr  = flag.String("red", "./data/renamed", "renamed entities json data directory")
// 		dirRnColPtr  = flag.String("rcd", "./data/renamed/collections", "renamed collections json data directory")
// 	)

// 	flag.Parse()

// 	//////////////////////////////////////////////////////////////

// 	dirOriEnt, dirRnEnt := *dirOriEntPtr, *dirRnEntPtr

// 	gio.MustCreateDir(dirRnEnt)

// 	// clear destination dir for putting renamed file
// 	lk.FailOnErr("%v", fd.RmFilesIn(dirRnEnt, false, true, "json"))

// 	// make sure each file's name is its entity value
// 	FixFileName(dirOriEnt, dirRnEnt)

// 	//////////////////////////////////////////////////////////////

// 	dirOriCol, dirRnCol := *dirOriColPtr, *dirRnColPtr

// 	gio.MustCreateDir(dirRnCol)

// 	// clear destination dir for putting renamed file
// 	lk.FailOnErr("%v", fd.RmFilesIn(dirRnCol, false, true, "json"))

// 	// make sure each file's name is its entity value
// 	FixFileName(dirOriCol, dirRnCol)

// 	// <------------------------------------------------------------------------------------->

// 	mChk := map[string][]string{
// 		dirRnEnt: {"Element", "Object", "Abstract Element"},
// 		dirRnCol: {"Collection"},
// 	}

// 	for _, dir := range []string{dirRnEnt, dirRnCol} {
// 		fs, err := os.ReadDir(dir)
// 		lk.FailOnErr("%v", err)
// 		for _, f := range fs {
// 			if fname := f.Name(); strings.HasSuffix(fname, ".json") {
// 				fpath := filepath.Join(dir, fname)
// 				data, err := os.ReadFile(fpath)
// 				lk.FailOnErr("%v", err)
// 				lk.WarnOnErrWhen(NotIn(gjson.Get(string(data), "Metadata.Type").String(), mChk[dir]...), "%v@%s", errors.New("ERROR TYPE"), fpath)
// 			}
// 		}
// 	}

// }

//////////////////////////////////////////////////////////////////////////////////////

//
// 2) renamed => out / err --- cmd: `go build -o preproc`
//
func main() {

	var (
		dirInEntPtr  = flag.String("ie", "./data/renamed", "input entities data directory")
		dirInColPtr  = flag.String("ic", "./data/renamed/collections", "input collections data directory")
		dirOutEntPtr = flag.String("oe", "./data/out", "output entities data directory")
		dirOutColPtr = flag.String("oc", "./data/out/collections", "output collections data directory")
		dirErrEntPtr = flag.String("ee", "./data/err", "error entities data directory")
		dirErrColPtr = flag.String("ec", "./data/err/collections", "error collections data directory")
	)

	flag.Parse()

	mInOut := map[string]string{
		*dirInEntPtr: *dirOutEntPtr,
		*dirInColPtr: *dirOutColPtr,
	}
	mInErr := map[string]string{
		*dirInEntPtr: *dirErrEntPtr,
		*dirInColPtr: *dirErrColPtr,
	}

	gio.MustCreateDirs(*dirOutEntPtr, *dirOutColPtr, *dirErrEntPtr, *dirErrColPtr)

	for I, dir := range []string{*dirInEntPtr, *dirInColPtr} {

		out := mInOut[dir]       // "out" is final output directory for ingestion
		errfolder := mInErr[dir] // "err" is for incorrect format json dump into

		lk.FailOnErr("%v", fd.RmFilesIn(out, false, false))
		lk.FailOnErr("%v", fd.RmFilesIn(errfolder, false, false))

		Preproc(dir, out, errfolder)

		/////////////////////////////////////////////////////////////////////

		DumpClassLinkage(out, "class-link.json")

		/////////////////////////////////////////////////////////////////////

		DumpPathValue(out, "path_val")

		/////////////////////////////////////////////////////////////////////

		if I == 0 {
			DumpCollection(out, "collection-entities.json")
		}
	}
}
