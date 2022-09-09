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
            { Root: "go.h4n.io/zetman", Type: "git", Source: "github.com/hbjydev/zetman" },
            { Root: "go.h4n.io/centra/component-base", Type: "git", Source: "github.com/centra-oss/component-base" },
        },
    }

    m := http.NewServeMux()

    tmpl, err := template.New("go-import").Parse(html)
    if err != nil {
        log.Fatal(err)
        return
    }

    for _, p := range c.GoPaths {
        parts := strings.Split(p.Root, "/")
        routeParts := parts[1:]
        route := fmt.Sprintf("/%v", strings.Join(routeParts, "/"))
        log.Printf("registering handler for: %v", route)

        handle := func (res http.ResponseWriter, req *http.Request) {
            q := req.URL.Query()
            if !q.Has("go-get") {
                http.Redirect(res, req, fmt.Sprintf("https://pkg.go.dev/%v", p.Root), http.StatusTemporaryRedirect)
                return
            }

            tmpl.Execute(res, p)
        }

        m.HandleFunc(route, handle)
        m.HandleFunc(fmt.Sprintf("%v/", route), handle)
    }

    s := http.Server{
        Addr: "0.0.0.0:8080",
        Handler: m,
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

