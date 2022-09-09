package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"
)

const version = "v0.0.0-DEV"

type Config struct {
    GoPaths []GoPath `json:"goPaths" yaml:"goPaths"`
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

    c := Config{
        GoPaths: []GoPath{
            { Root: "go.h4n.io/zetman", Type: "git", Source: "https://github.com/hbjydev/zetman" },
            { Root: "go.h4n.io/centra/component-base", Type: "git", Source: "https://github.com/centra-oss/component-base" },
        },
    }

    // m := http.NewServeMux()

    tmpl, err := template.New("go-import").Parse(html)
    if err != nil {
        log.Fatal(err)
        return
    }

    httpHandler := http.HandlerFunc(func (w http.ResponseWriter, r *http.Request) {
        url := r.URL.Path
        log.Printf("GET %v", url)

        var goPath GoPath

        for _, p := range c.GoPaths {
            parts := strings.Split(p.Root, "/")
            routeParts := parts[1:]
            route := fmt.Sprintf("/%v", strings.Join(routeParts, "/"))

            log.Printf("logging test for %v", route)

            if strings.HasPrefix(url, route) {
                goPath = p
            }
        }

        q := r.URL.Query()
        if !q.Has("go-get") {
            http.Redirect(w, r, fmt.Sprintf("https://pkg.go.dev/%v", goPath.Root), http.StatusTemporaryRedirect)
            return
        }

        tmpl.Execute(w, goPath)
    })

    s := http.Server{
        Addr: "0.0.0.0:8080",
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

