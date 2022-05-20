package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	fd "github.com/digisan/gotk/filedir"
	gio "github.com/digisan/gotk/io"
	"github.com/digisan/gotk/strs"
	jt "github.com/digisan/json-tool"
	lk "github.com/digisan/logkit"
	"github.com/tidwall/gjson"
)

// make sure each file's name is its entity value
func FixFilename(datadir, odir string) {

	files, err := os.ReadDir(datadir)
	if err != nil {
		log.Fatalln(err)
	}

	for _, file := range files {
		fpath := filepath.Join(datadir, file.Name())
		// fmt.Println(fpath)

		if strings.HasSuffix(fpath, ".json") {
			data, err := os.ReadFile(fpath)
			if err != nil {
				log.Fatalln(err)
			}
			lk.Log("reading...  %s", fpath)

			entity := gjson.Get(string(data), "Entity").String()
			fname := entity + ".json"
			if len(odir) == 0 {
				odir = strs.SplitPartFromLast(fpath, "/", 2)
			}
			fpathNew := filepath.Join(odir, fname)
			lk.Log("destination...  %s", fpathNew)

			lk.FailOnErrWhen(fd.FileExists(fpathNew), "%v", fmt.Errorf("[%s] is already existing", fpathNew))

			// copy
			if err = os.WriteFile(fpathNew, data, os.ModePerm); err != nil {
				log.Fatalln(err)
			}

			// move
			// if err = os.Rename(fpath, fpathNew); err != nil {
			// 	log.Fatalln(err)
			// }
		}
	}
}

// 1) Escape quotation marks used around HTML attributes like so <img src=\"someimage.png\" />

// 2) Escape the forward slash in HTML end tags. <div>Hello
//    World!<\/div>. This is an ancient artifact of an old HTML spec that didn't want HTML parsers to get confused when putting strings in a <SCRIPT> tag. For some reason, today’s browsers still like it.

// 3) This one was totally bizarre. You should include a space between the tag name and the slash on self-closing tags. I have no idea why this is, but on MOST modern browsers, if you try using javascript to append a <li> tag as a child of an unordered list that is formatted like so: <ul/>, it won't work. It gets added to the DOM after the ul tag. But, if the code looks like this: <ul /> (notice the space before the /), everything works fine. Very strange indeed.

// 4) Be sure to encode any quotation marks that might be included in (bad) HTML content. This is the only thing that would really break the JSON by accidentally terminating the string early. Any " characters should be encoded as &quot; if it is meant to be included as HTML content.

func rmLF(data []byte) []byte {
	data = bytes.ReplaceAll(data, []byte{'\n'}, []byte{})
	data = bytes.ReplaceAll(data, []byte{'\r'}, []byte{})
	return data
}

func escQuInHTML(ori string) string {
	r := regexp.MustCompile(`<[\w\d]+\s[^>]+>`)
	return r.ReplaceAllStringFunc(ori, func(s string) string {
		// fmt.Println("---", s)
		s = strings.ReplaceAll(s, `"`, `\"`)
		s = strings.ReplaceAll(s, `\\`, `\`)
		return s
	})
}

func fixErrComma(s string) string {
	r := regexp.MustCompile(`,\s*[\}\]]`)
	spanList := r.FindAllStringIndex(s, -1)
	// for _, span := range spanList {
	// 	b, e := span[0], span[1]
	// 	fmt.Println(s[b:e])
	// }
	spanls := [][2]int{}
	for _, span := range spanList {
		spanls = append(spanls, [2]int{span[0], span[1] - 1})
	}
	return strs.RangeReplace(s, spanls, []string{" "})
}

func rmPtag(ori string) string {
	r := regexp.MustCompile(`</p>\s*<p>`)
	ori = r.ReplaceAllStringFunc(ori, func(s string) string {
		return "<br>"
	})
	r = regexp.MustCompile(`<p>`)
	ori = r.ReplaceAllStringFunc(ori, func(s string) string {
		return ""
	})
	r = regexp.MustCompile(`</p>`)
	ori = r.ReplaceAllStringFunc(ori, func(s string) string {
		return ""
	})
	return ori
}

func Preproc(datadir, odir, edir string) error {

	files, err := os.ReadDir(datadir)
	if err != nil {
		return err
	}

	for _, f := range files {
		if fpath := filepath.Join(datadir, f.Name()); strings.HasSuffix(fpath, ".json") {
			lk.Log("processing...  %v", fpath)

			data, err := os.ReadFile(fpath)
			if err != nil {
				return err
			}
			if len(data) == 0 {
				continue
			}

			data = rmLF(data)
			data = []byte(rmPtag(string(data)))
			data = []byte(escQuInHTML(string(data)))
			data = []byte(fixErrComma(string(data)))

			if !jt.IsValid(data) {
				gio.MustCreateDir(edir)
				outname := filepath.Base(fpath)
				out := filepath.Join(edir, outname)
				os.WriteFile(out, data, os.ModePerm)
				lk.FailOnErr("%v", fmt.Errorf("json error@ %s", fpath))
			}

			// save
			gio.MustCreateDir(odir)
			outname := filepath.Base(fpath)
			out := filepath.Join(odir, outname)
			os.WriteFile(out, data, os.ModePerm)

			lk.Log("%s is processed & stored", out)
		}
	}
	return nil
}
