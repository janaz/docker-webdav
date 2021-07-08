package main

import (
	"log"
	"net/http"
	"os"

	"golang.org/x/net/webdav"
)

type methodMux map[string]http.Handler

func (m *methodMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if h, ok := (*m)[r.Method]; ok {
		h.ServeHTTP(w, r)
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func main() {
	listen := os.Getenv("LISTEN")
	root := os.Getenv("ROOT")
	prefix := os.Getenv("PREFIX")

	files := http.StripPrefix(prefix, http.FileServer(http.Dir(root)))
	webdav := &webdav.Handler{
		Prefix:     prefix,
		FileSystem: webdav.Dir(root),
		LockSystem: webdav.NewMemLS(),
		Logger: func(r *http.Request, err error) {
			if err != nil {
				log.Printf("r=%v err=%v", r, err)
			}
		},
	}
	loggedWebdav := logRequestHandler(webdav)
	mux := methodMux(map[string]http.Handler{
		"GET":       logRequestHandler(files),
		"OPTIONS":   loggedWebdav,
		"PROPFIND":  loggedWebdav,
		"PROPPATCH": loggedWebdav,
		"MKCOL":     loggedWebdav,
		"COPY":      loggedWebdav,
		"MOVE":      loggedWebdav,
		"LOCK":      loggedWebdav,
		"UNLOCK":    loggedWebdav,
		"DELETE":    loggedWebdav,
		"PUT":       loggedWebdav,
	})

	if err := http.ListenAndServe(listen, &mux); err != nil {
		log.Fatal(err)
	}
}

func logRequestHandler(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {

		// call the original http.Handler we're wrapping
		h.ServeHTTP(w, r)

		// gather information about request and log it
		uri := r.URL.String()
		method := r.Method
		// ... more information
		log.Printf("%s %s", method, uri)
	}

	// http.HandlerFunc wraps a function so that it
	// implements http.Handler interface
	return http.HandlerFunc(fn)
}
