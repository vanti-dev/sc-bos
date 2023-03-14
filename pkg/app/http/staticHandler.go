package http

import (
	"net/http"
	"os"
	"path"
	"strings"
)

func NewStaticHandler(staticPath string) http.Handler {
	return SPAFileServer(http.Dir(staticPath), func(w http.ResponseWriter, r *http.Request) bool {
		r.URL.Path = "/"
		return true
	})
}

// SPAHandler provides the function signature for passing to the SPAFileServer
type SPAHandler = func(w http.ResponseWriter, r *http.Request) (doDefaultFileServe bool)

// SPAFileServer wraps the http.FileServer checking to see if the url path exists first.
//
// # If the file fails to exist it calls the supplied handlerSPA function
//
// The implementation can choose to either modify the request, e.g. change the URL path and return true to have the
// default FileServer handling to still take place, or return false to stop further processing, for example if you
// wanted to write a custom response
//
// This code borrows heavily from https://gist.github.com/lummie/91cd1c18b2e32fa9f316862221a6fd5c
func SPAFileServer(root http.FileSystem, handlerSPA SPAHandler) http.Handler {
	fs := http.FileServer(root)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//make sure the url path starts with /
		upath := r.URL.Path
		if !strings.HasPrefix(upath, "/") {
			upath = "/" + upath
			r.URL.Path = upath
		}
		upath = path.Clean(upath)

		// attempt to open the file via the http.FileSystem
		f, err := root.Open(upath)
		if err != nil {
			if os.IsNotExist(err) {
				// call handler
				if handlerSPA != nil {
					doDefault := handlerSPA(w, r)
					if !doDefault {
						return
					}
				}
			}
		}

		// close if successfully opened
		if err == nil {
			f.Close()
		}

		// default serve
		fs.ServeHTTP(w, r)
	})
}
