package devices

import (
	"fmt"
	"strconv"
	"strings"
)

// dePath represents a deconstructed path used to extract an array Index if present
type dePath struct {
	Before string
	After  string
	Found  bool
	Index  int
	Next   string
}

// depath deconstructs a path to extract an array Index if present
// returns true for Found if a possible integer is Found wrapped in [ ]
func depath(path string) dePath {
	before, after, _ := strings.Cut(path, ".")

	indexStr := strings.Builder{}
	buildIndex := false
	foundIndex := false
	outPath := strings.Builder{}
	addToOutPath := true
	indexOfNum := -1
	for idx, char := range before {
		if char == '[' {
			buildIndex = true
			addToOutPath = false
			indexOfNum = idx + 1
		} else if buildIndex {
			if char == '-' || (char >= '0' && char <= '9') {
				indexStr.WriteRune(char)
			} else if char == ']' {
				buildIndex = false
				foundIndex = true
				break
			} else {
				foundIndex = false
				buildIndex = false
			}
		}

		if addToOutPath {
			outPath.WriteRune(char)
		}
	}

	if !foundIndex {
		return dePath{
			Before: before,
			After:  after,
			Found:  false,
			Index:  -1,
			Next:   after,
		}
	}

	iStr := indexStr.String()
	index, err := strconv.Atoi(iStr)

	if err != nil || index < 0 {
		return dePath{
			Before: outPath.String(),
			After:  after,
			Found:  foundIndex,
			Index:  -1,
			Next:   after,
		}
	}

	if indexOfNum >= 1 {
		return dePath{
			Before: outPath.String(),
			After:  fmt.Sprintf("[%d].%s", index, after),
			Found:  foundIndex,
			Index:  index,
			Next:   after,
		}
	}

	if after != "" {
		return dePath{
			Before: outPath.String(),
			After:  after,
			Found:  foundIndex,
			Index:  -1,
			Next:   after,
		}
	}

	return dePath{
		Before: before,
		After:  after,
		Found:  foundIndex,
		Index:  -1,
		Next:   after,
	}
}
