package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"
)

const version = "v0.0.0-DEV"

var (
	addr *string = flag.String("http.listen-addr", "0.0.0.0:8080", "the address to listen on")
	host *string = flag.String("http.public-host", "go.h4n.io", "the public hostname of this nail-care instance")
)

type Config struct {
	PublicHostname string   `json:"publicHostname" yaml:"publicHostname"`
	GoPaths        []GoPath `json:"goPaths" yaml:"goPaths"`
}

type GoPath struct {
	// Root is the import root to advertise to Go.
	Root string `json:"root" yaml:"root"`

	// Type specifies the type of import (i.e. git)
	Type string `json:"type" yaml:"type"`

	// Source is the source URL of the path.
	Source string `json:"source" yaml:"source"`
}

func main() {
	log.Printf("nail_care version %s", version)
	flag.Parse()

    var paths []GoPath

	a := flag.Args()

    if len(a) == 0 {
        log.Fatal("no paths defined")
    }

    // Parse all command-line flags for config
	for _, v := range a {
        // Format: nail_care --http.listen=:8080 --http.host=go.h4n.io zetman,git,https://github.com
        // Argument format: [path],[type],[source]

        spl := strings.Split(v, ",")

        if len(spl) % 3 != 0 {
            log.Fatal("invalid path, needs three csv parts")
            return
        }

        lenCond := (len(spl[0]) == 0 || len(spl[1]) == 0 || len(spl[2]) == 0)
        if lenCond {
            log.Fatal("invalid path, must specify three csv parts")
        }

        path := GoPath{
            Root: spl[0],
            Type: spl[1],
            Source: spl[2],
        }

        paths = append(paths, path)
	}

	c := Config{
		PublicHostname: *host,
		GoPaths: paths,
	}

	tmpl, err := template.New("go-import").Parse(html)
	if err != nil {
		log.Fatal(err)
		return
	}

	httpHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		url := r.URL.Path
		log.Printf("GET %v", url)

		var goPath *GoPath

		// Find the go path for the current request
		for _, p := range c.GoPaths {
			route := fmt.Sprintf("/%v", p.Root)
			if strings.HasPrefix(url, route) {
				goPath = &p
			}
		}

		// If no path was found for that route, throw an error.
		if goPath == nil {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("No such package found."))
			return
		}

		// If the request is not from `go get`, redirect to pkg.go.dev.
		q := r.URL.Query()
		if !q.Has("go-get") {
			http.Redirect(
				w,
				r,
				fmt.Sprintf( // build pkg.go.dev path
					"https://pkg.go.dev/%s/%s",
					c.PublicHostname,
					goPath.Root,
				),
				http.StatusTemporaryRedirect, // not a permanent redirect
			)
			return
		}

		// Otherwise, return the required HTML for Go to install the packages.
		tmpl.Execute(w, goPath)
	})

	s := http.Server{
		Addr:    *addr,
		Handler: httpHandler,
	}

	log.Fatal(s.ListenAndServe())
}

const html = `
<!DOCTYPE html>
<html>
    <head>
        <meta http-equiv="Content-Type" content="text/html; charset=utf-8"/>
        <meta name="go-import" content="{{.Root}} {{.Type}} {{.Source}}">
    </head>
    <body></body>
</html>
`
