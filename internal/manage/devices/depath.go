package devices

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var arrIndexRegx = regexp.MustCompile("\\[(-?[0-9]+)]")

// dePath represents a deconstructed path used to extract an array Index if present
type dePath struct {
	Before string
	After  string
	Found  bool
	Index  int
}

// depath deconstructs a path to extract an array Index if present
// returns true for Found if a possible integer is Found wrapped in [ ]
func depath(path string) dePath {
	before, after, _ := strings.Cut(path, ".")

	if arrIndexRegx.MatchString(path) {
		matches := arrIndexRegx.FindStringSubmatch(path)

		if len(matches) < 2 {
			return dePath{
				Before: before,
				After:  after,
				Found:  true,
				Index:  -1,
			}
		}

		index, err := strconv.ParseInt(matches[1], 10, 32)
		if err == nil && index > -1 {
			matchedIndices := arrIndexRegx.FindStringIndex(path)

			if matchedIndices == nil || len(matchedIndices) < 2 {
				return dePath{
					Before: before,
					After:  after,
					Found:  true,
					Index:  -1,
				}
			}
			if matchedIndices[0] == 0 {
				// An Index is Found at the start of path
				// return Before,After only
				return dePath{
					Before: before,
					After:  after,
					Found:  true,
					Index:  int(index),
				}
			}
			return dePath{
				Before: arrIndexRegx.ReplaceAllString(before, ""),
				After:  fmt.Sprintf("[%d].%s", index, after),
				Found:  true,
				Index:  int(index),
			}
		}

		return dePath{
			Before: arrIndexRegx.ReplaceAllString(before, ""),
			After:  after,
			Found:  true,
			Index:  -1,
		}
	}

	return dePath{
		Before: before,
		After:  after,
		Found:  false,
		Index:  -1,
	}
}
