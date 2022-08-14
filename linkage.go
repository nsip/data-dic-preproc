package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	. "github.com/digisan/go-generics/v2"
	fd "github.com/digisan/gotk/filedir"
	lk "github.com/digisan/logkit"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

func SupClsCol(js string) map[string][]string {
	ret := make(map[string][]string)
	entity := gjson.Get(js, "Entity").String()
	supers := gjson.Get(js, "Metadata.Superclass").Array()
	for _, super := range supers {
		ret[entity] = append(ret[entity], super.String())
	}
	return ret
}

func SwapES(m map[string][]string) map[string][]string {
	ret := make(map[string][]string)
	for entity, supers := range m {
		for _, super := range supers {
			ret[super] = append(ret[super], entity)
		}
	}
	return ret
}

func EntitiesDescArr(fpaths ...string) []map[string][]string {
	mEDs := make([]map[string][]string, 0, len(fpaths))
	for _, path := range fpaths {
		data, err := os.ReadFile(path)
		lk.FailOnErr("%v", err)
		js := string(data)
		mES := SupClsCol(js)
		mEDs = append(mEDs, SwapES(mES))
	}
	return mEDs
}

func EntityDesc(fpaths ...string) map[string][]string {

	mEDs := EntitiesDescArr(fpaths...)
	keys := []string{}
	for _, mED := range mEDs {
		for k := range mED {
			keys = append(keys, k)
		}
	}
	keys = Settify(keys...)

	mEDsKey := make(map[string][]string)
	for _, key := range keys {
		for _, mED := range mEDs {
			for k := range mED {
				if k == key {
					mEDsKey = MapMergeOnValSlc(mEDsKey, mED)
				}
			}
		}
	}

	/// for testing deeply
	// mEDsKey["Campus"] = []string{"Sydenham Campus", "Hillside Campus", "Taylors Camplus"}
	// mEDsKey["Sydenham Campus"] = []string{"Sydenham-Hillside Campus 1"}
	// mEDsKey["Hillside Campus"] = []string{"Sydenham-Hillside Campus 2"}
	///

	return mEDsKey
}

type List []string

func (ls List) String() string {
	sb := strings.Builder{}
	for i, ele := range ls {
		sb.WriteString(ele)
		if i < len(ls)-1 {
			sb.WriteString("--")
		}
	}
	return sb.String()
}

func LinkEntity(mED map[string][]string, keyEntity string, ancestry List, linkCol *List) {
	lookfor := keyEntity
	for entity, descList := range mED {
		if entity == lookfor {
			for _, desc := range descList {
				link := fmt.Sprintf("%s--%s--%s", ancestry, entity, desc)
				link = strings.TrimLeft(link, "--")
				// fmt.Println(link)
				*linkCol = append(*linkCol, link)
				lookfor = desc
				delete(mED, entity)
				ancestry = append(ancestry, keyEntity)
				LinkEntity(mED, lookfor, ancestry, linkCol)
				ancestry = ancestry[:len(ancestry)-1]
			}
		}
	}
}

func RmPartialLink(linkCol []string) []string {
AGAIN:
	for _, linkCheck := range linkCol {
		for _, linkCompare := range linkCol {
			if linkCheck != linkCompare {
				if strings.HasPrefix(linkCompare, linkCheck) ||
					strings.HasSuffix(linkCompare, linkCheck) ||
					strings.Contains(linkCompare, "--"+linkCheck+"--") {
					DelOneEle(&linkCol, linkCheck)
					goto AGAIN
				}
			}
		}
	}
	return linkCol
}

func LinkEntities(fpaths ...string) (out []string) {
	mED := EntityDesc(fpaths...)
	for k := range mED {
		linkCol := &List{}
		LinkEntity(MapCopy(mED), k, List{}, linkCol)
		out = append(out, RmPartialLink(*linkCol)...)
	}
	return RmPartialLink(out)
}

func TrimEntityPaths(mEntityPaths map[string][]string) map[string][]string {
	for k, v := range mEntityPaths {
		for i := 0; i < len(v); i++ {
			p := v[i]
			if strings.HasPrefix(p, k+"--") {
				v[i] = k
			}
			if pos := strings.Index(p, "--"+k+"--"); pos >= 0 {
				v[i] = p[:pos+len(k)+2]
			}
		}
		mEntityPaths[k] = Settify(v...)
	}
	return mEntityPaths
}

type Node struct {
	Branch   string
	Children []string
}

func CleanUpEntityPaths(mEntityPaths map[string][]string) map[string]Node {
	m := make(map[string]Node)
	for entity, paths := range mEntityPaths {
		// one node
		node := Node{}
		for _, path := range paths {
			ss := strings.Split(path, "--")
			for i, s := range ss {
				if s == entity {
					if i == len(ss)-1 {
						node.Branch = path
						node.Children = []string{}
					} else {
						node.Branch = strings.Join(ss[:i+1], "--")
						// child := strings.ReplaceAll(ss[i+1], ".", "[dot]")
						child := ss[i+1]
						node.Children = append(node.Children, child)
					}
				}
			}
		}
		node.Children = Settify(node.Children...)
		m[entity] = node
	}
	return m
}

func Link2JSON(linkCol []string, path string) (out string, err error) {

	mEntityPathsCol := []map[string]string{}

	for _, link := range linkCol {
		for _, entity := range strings.Split(link, "--") {
			// fmt.Println(entity, ":", link)
			m := make(map[string]string)
			m[entity] = link
			mEntityPathsCol = append(mEntityPathsCol, m)
		}
	}

	mEntityPaths := MapMerge(mEntityPathsCol...)

	// {
	// 	// only keep value as tree branch list to terminate at current 'key' node.
	// 	mEntityPaths = TrimEntityPaths(mEntityPaths)

	// 	for k, v := range mEntityPaths {
	// 		fmt.Println(k, v)
	// 		fmt.Println()

	// 		if strings.Contains(k, ".") {
	// 			k = strings.ReplaceAll(k, ".", "[dot]")
	// 		}

	// 		out, err = sjson.Set(out, k, v)
	// 		lk.FailOnErr("%v", err)
	// 	}
	// }

	{
		mEntityNode := CleanUpEntityPaths(mEntityPaths)

		for entity, node := range mEntityNode {
			// fmt.Println(entity, node)
			// fmt.Println()

			if strings.Contains(entity, ".") {
				entity = strings.ReplaceAll(entity, ".", "[dot]")
			}

			out, err = sjson.Set(out, entity, node)
			lk.FailOnErr("%v", err)
		}
	}

	return out, nil // change "." to "[dot]" from each key, otherwise, mongodb stores unexpected...
}

func DumpClassLinkage(idir, ofname string) {
	files, _, err := fd.WalkFileDir(idir, false)
	lk.FailOnErr("%v", err)
	linkCol := LinkEntities(files...)
	js, err := Link2JSON(linkCol, "")
	lk.FailOnErr("%v", err)
	lk.FailOnErr("%v", os.WriteFile(filepath.Join(idir, ofname), []byte(js), os.ModePerm))
}
