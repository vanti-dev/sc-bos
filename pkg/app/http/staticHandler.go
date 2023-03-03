package http

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

// staticHandler implements the http.Handler interface, so we can use it
// to respond to HTTP requests. The path to the static directory and
// path to the index file within that static directory are used to
// serve the SPA in the given static directory.
type staticHandler struct {
	staticPath string
	indexPath  string
	logger     *zap.Logger
}

func ServeStaticSite(c StaticHostingConfig, mux *mux.Router, logger *zap.Logger) {
	logger.Info("Serving static site", zap.Any("config", c))
	static := NewStaticHandler(c.FilePath, logger)
	mux.PathPrefix(c.Path).Handler(http.StripPrefix(c.Path, static))
}

func NewStaticHandler(staticPath string, logger *zap.Logger) *staticHandler {
	return &staticHandler{
		staticPath,
		"index.html",
		logger,
	}
}

// ServeHTTP inspects the URL path to locate a file within the static dir
// on the SPA handler. If a file is found, it will be served. If not, the
// file located at the index path on the SPA handler will be served. This
// is suitable behavior for serving an SPA (single page application).
func (h staticHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// get the absolute path to prevent directory traversal
	p := filepath.Clean(r.URL.Path)

	// prepend the path with the path to the static directory
	path := filepath.Join(h.staticPath, p)

	h.logger.Debug("looking for file", zap.String("req", r.URL.Path), zap.String("path", path))

	// check whether a file exists at the given path
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		h.logger.Debug("File not found, serving index")
		// file does not exist, serve index.html
		http.ServeFile(w, r, filepath.Join(h.staticPath, h.indexPath))
		return
	} else if err != nil {
		// if we got an error (that wasn't that the file doesn't exist) stating the
		// file, return a 500 internal server error and stop
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// otherwise, use http.FileServer to serve the static dir
	http.FileServer(http.Dir(h.staticPath)).ServeHTTP(w, r)
}
