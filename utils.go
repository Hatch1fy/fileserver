package fileserver

import (
	"mime"
	"net/http"
	"path/filepath"
	"strings"
)

// canAccept returns if a request supports gzip and/or webp responses
func canAccept(req *http.Request) (gzip, webp bool) {
	gzip = canAcceptGZip(req)
	webp = canAcceptWebP(req)
	return
}

// canAcceptWebP returns if a request supports webp responses
func canAcceptWebP(req *http.Request) (ok bool) {
	var header []string
	// Get accept header
	if header, ok = req.Header["Accept"]; !ok {
		// Accept header does not exist, return
		return
	}

	for _, accept := range header {
		if strings.Index(accept, "image/webp") > -1 {
			// This entry is image/webp, return
			return true
		}
	}

	return false
}

// canAcceptGZip returns if a request supports gzip responses
func canAcceptGZip(req *http.Request) (ok bool) {
	var header []string
	// Get accept header
	if header, ok = req.Header["Accept-Encoding"]; !ok {
		// Accept header does not exist, return
		return
	}

	for _, encoding := range header {
		if encoding == "gzip" {
			// This entry is not image/webp, continue
			return true
		}
	}

	return false
}

// isCachable returns if a given extension is cachable
func isCachable(ext string) (cachable bool) {
	switch ext {
	// Block of cachable extensions
	case ".js":
	case ".css":
	case ".png":
	case ".jpg":
	case ".webp":
	case ".ttf":
	case ".wof":

	// Extension did not match cachable block
	default:
		return false
	}

	return true
}

func getKey(key, option string) string {
	if strings.Index(option, "webp") > -1 {
		// Option has webp, strip the original extension
		key = stripExt(key)
	}

	return key + option
}

func getOptions(req *http.Request) (options []string) {
	gzip, webp := canAccept(req)
	// Switch on acceptable responses
	switch {
	case !gzip && !webp:
		options = noOptions
	case gzip && !webp:
		options = gzipOptions
	case !gzip && webp:
		options = webpOptions
	case gzip && webp:
		options = webpGZipOptions
	}

	return
}

func setHeaders(w http.ResponseWriter, filename, ext string) {
	header := w.Header()
	if ext == ".gz" {
		// We are serving a gzipped file, set content encoding
		header.Set("Content-Encoding", "gzip")
		// Strip gzip extension so MIME type can be properly determined
		ext = filepath.Ext(stripExt(filename))
	}

	// Set the content type of the file we're serving
	header.Set("Content-Type", mime.TypeByExtension(ext))

	if isCachable(ext) {
		// Extension is cachable, set cache control
		header.Set("Cache-Control", "public,max-age=3600")
	}
}

// stripExt will strip the extension from a file
func stripExt(key string) string {
	// Get key's extension
	ext := filepath.Ext(key)
	// Return key with extension removed
	return key[:len(key)-len(ext)]
}
