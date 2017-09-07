package atc

import (
	"errors"
	"github.com/adolphlxm/atc/context"
	"net/http"
	"os"
	"path"
	"strings"
)

var errNotStaticRequest = errors.New("request not a static file request")

// frontStaticRouter is the default static file handler - this is the last line of handlers
func frontStaticRouter(c *context.Context) error {
	if c.Method() != "GET" && c.Method() != "HEAD" {
		return errNotStaticRequest
	}

	var isFront bool

	requestPath := path.Clean(c.Path())
	localPath := "./front" + requestPath
	// special processing : favicon.ico/robots.txt  can be in any static dir
	if requestPath == "/favicon.ico" || requestPath == "/robots.txt" {
		serveFile(localPath, c)
		return nil
	}

	//Matching static file
	for _, prefix := range Aconfig.FrontDir {
		if strings.HasPrefix(requestPath, "/"+prefix) {
			isFront = true
		}
	}

	if isFront {
		return serveFile(localPath, c)

	}

	return errNotStaticRequest
}

func serveFile(localPath string, c *context.Context) error {
	f, err := os.Stat(localPath)

	if f == nil {
		htmlLocalPath := strings.Replace(localPath, "/", "", 0) + "." + Aconfig.FrontSuffix
		if f, _ = os.Stat(htmlLocalPath); f != nil {
			// If the file exists and we can access it, serve it
			http.ServeFile(c.ResponseWriter, c.Request, htmlLocalPath)
			return nil
		}

		if os.IsNotExist(err) {
			//http.NotFound(c.ResponseWriter(), c.Request())
		}
		return errNotStaticRequest
	}

	if f.IsDir() {
		if f, _ = os.Stat(path.Join(localPath, "index.html")); f == nil && !Aconfig.FrontDirectory {
			http.NotFound(c.ResponseWriter, c.Request)
			return errNotStaticRequest
		}

	}

	// If the file exists and we can access it, serve it
	http.ServeFile(c.ResponseWriter, c.Request, localPath)
	return nil
}
