// Application vet-config checks Smart Core config files for common errors.
// In general the tool will recursively scan a directory for `.json` files and analyse the contents looking for potential problems.
//
// # Checks
//
// Duplicates:
//
// The tool finds arrays that look like sets (have some id/key field) and checks for duplicate items.
// What the tool defines as keys can be configured using the `-key` flag, which can be used multiple times.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"

	"github.com/qri-io/jsonpointer"
	"golang.org/x/exp/maps"
)

var (
	dir    = flag.String("dir", ".", "directory to scan")
	suffix = flag.String("suffix", ".json", "file suffix to scan")
)

var (
	ignore stringList
	keys   = stringList{"{name,kind}", "id", "name", "host", "test"}
)

func init() {
	flag.Var(&ignore, "ignore", "ignore glob paths")
	flag.Var(&keys, "key", "key fields to look for")
}

func main() {
	flag.Parse()
	// scan through all files in the directory
	err := filepath.WalkDir(*dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		localPath, err := filepath.Rel(*dir, path)
		if err != nil {
			return err
		}
		for _, ig := range ignore {
			if match, err := filepath.Match(ig, localPath); err != nil {
				return err
			} else if match {
				return filepath.SkipDir
			}
		}
		if d.IsDir() {
			return nil
		}
		if !strings.HasSuffix(path, *suffix) {
			return nil
		}

		in, err := os.Open(path)
		if err != nil {
			return err
		}
		defer in.Close()

		var data any
		err = json.NewDecoder(in).Decode(&data)
		if err != nil {
			return err
		}

		inspect(localPath, jsonpointer.NewPointer(), data)
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}

	var setDupes []dupe
	for k, v := range results {
		if len(v.SetKeyDupes) > 0 {
			setDupes = append(setDupes, dupe{k, v.SetKey, v.SetKeyDupes})
		}
	}

	// arrange dupes by file
	type dupeByFile struct {
		file  string
		dupes map[string][]string // keyed by path/key, values are the duplicated values
	}
	var dupesByFile []dupeByFile // sorted by .file
	for _, setDupe := range setDupes {
		for f, dupes := range setDupe.dupesByFile {
			dupeProp := path.Join(setDupe.path, "*", setDupe.key)
			i, found := sort.Find(len(dupesByFile), func(i int) int {
				return strings.Compare(f, dupesByFile[i].file)
			})
			if !found {
				dupesByFile = append(dupesByFile, dupeByFile{})
				copy(dupesByFile[i+1:], dupesByFile[i:])
				dupesByFile[i] = dupeByFile{f, map[string][]string{}}
			}
			for _, id := range dupes {
				dupesByFile[i].dupes[dupeProp] = addString(dupesByFile[i].dupes[dupeProp], id)
			}
		}
	}

	// print out the results
	for _, fileDupes := range dupesByFile {
		fmt.Printf("dupes in %q:\n", fileDupes.file)
		keys := maps.Keys(fileDupes.dupes)
		sort.Strings(keys)
		for _, k := range keys {
			ids := fileDupes.dupes[k]
			for _, id := range ids {
				fmt.Printf("\t%q: %q\n", k, id)
			}
		}
	}
}

func inspect(file string, path jsonpointer.Pointer, data any) {
	switch data := data.(type) {
	case []any:
		if len(data) == 0 {
			return
		}
		// if this looks like an array of objects, check if it's a set by looking for a key
		if firstItem, ok := data[0].(map[string]any); ok {
			// find the key used for items in the set
			var setKey string
			for _, key := range keys {
				if _, ok := readKey(firstItem, key); ok {
					// check if the value is a string
					setKey = key
					break
				}
			}
			if setKey != "" {
				r := results.arraySet(path, setKey, len(data)).addFile(file)
				// check if the key value is duplicated in the set
				// Note: there's no guarantee that all items in the array are maps,
				// and no guarantee that all key property values are strings
				seen := map[string]struct{}{}
				for _, item := range data {
					m, ok := item.(map[string]any)
					if !ok {
						// not all items are objects
						r.MixedType = true
						break
					}
					s, ok := readKey(m, setKey)
					if !ok {
						// not all items have the key
						r.MixedType = true
						break
					}
					if _, ok := seen[s]; ok {
						// key value is duplicated
						r.addSetKeyDupe(file, s)
						break
					}
					seen[s] = struct{}{}
				}
			}
		}
		for _, v := range data {
			inspect(file, path.RawDescendant("*"), v)
		}
	case map[string]any:
		for k, v := range data {
			inspect(file, path.RawDescendant(k), v)
		}
	default:
		return
	}
}

type records map[string]*record

var results = records{}

type record struct {
	IsSet       bool
	SetKey      string
	SetKeyDupes map[string][]string

	MixedType bool // items in these arrays are not all the same type

	Count int // how many times this json path was found
	Items int // a sum of the length of all arrays found at this path

	Files []string // files where this record was found
}

func (r *record) addFile(file string) *record {
	r.Files = addString(r.Files, file)
	return r
}

func (r *record) addSetKeyDupe(file, val string) *record {
	if r.SetKeyDupes == nil {
		r.SetKeyDupes = map[string][]string{}
	}
	r.SetKeyDupes[file] = addString(r.SetKeyDupes[file], val)
	return r
}

type dupe struct {
	path        string
	key         string
	dupesByFile map[string][]string
}

func addString(ss []string, s string) []string {
	i, found := sort.Find(len(ss), func(i int) int {
		return strings.Compare(ss[i], s)
	})
	if found {
		return ss
	}
	ss = append(ss, "")
	copy(ss[i+1:], ss[i:])
	ss[i] = s
	return ss
}

func (r records) arraySet(path jsonpointer.Pointer, key string, l int) *record {
	rec, ok := r[path.String()]
	if !ok {
		rec = &record{}
		r[path.String()] = rec
	}
	rec.IsSet = true
	rec.SetKey = key
	rec.Items += l
	rec.Count++
	return rec
}

func readKey(d map[string]any, k string) (string, bool) {
	if k[0] == '{' {
		// composite key
		ks := strings.Split(k[1:len(k)-1], ",")
		var vs []string
		for _, k := range ks {
			v, ok := readKey(d, k)
			if !ok {
				return "", false
			}
			vs = append(vs, v)
		}
		return strings.Join(vs, ","), true
	}
	v, ok := d[k]
	if !ok {
		return "", false // doesn't exist
	}
	s, ok := v.(string)
	if !ok {
		return "", false // not a string
	}
	return s, true
}
