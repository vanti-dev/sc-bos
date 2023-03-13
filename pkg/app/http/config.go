package http

type StaticHostingConfig struct {
	// the location of the static files
	FilePath string `json:"filepath"`
	// the url path the files should be served on (e.g. "/my-site")
	Path string `json:"path"`
}

type CorsConfig struct {
	// DebugMode turns cors debug logging on/off
	DebugMode bool `json:"debugMode,omitempty"`
	// CorsOrigins specifies a list of allowed domains for passing to the CORS handler
	CorsOrigins []string `json:"corsOrigins,omitempty"`
}
