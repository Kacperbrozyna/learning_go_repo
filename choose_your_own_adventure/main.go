package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"

	"github.com/Kacperbrozyna/learning_go_repo/choose_your_own_adventure/cyoa"
)

const json_path = "gopher.json"
const useDefaultVals = true

const nonDefaultTemplate = `
<!DOCTYPE html>
<html>
    <head>
        <meta charset="utf-8">
        <title> Choose Your Own Adventure </title>
    </head>
    <body>
        <h1>{{.Title}}</h1>
        {{range .Paragraph}}
            <p>{{.}}</p>
        {{end}}
        <ul>    
        {{range .Options}}
            <li>
                <a href="/{{.Arc}}">{{.Text}}</a>
            </li>
        {{end}}
        </ul>
    </body>
</html>`

func nonDefaultPathFunc(r *http.Request) string {
	path := strings.TrimSpace(r.URL.Path)
	if path == "/story" || path == "/story/" {
		path = "story/intro"
	}

	return path[len("/story/"):]
}

func main() {

	port := flag.Int("port", 8080, "Port to start cyoa web app on")

	story, err := cyoa.ReadJsonStory(json_path)
	if err != nil {
		panic(err)
	}

	var handler http.Handler

	tpl := template.Must(template.New("").Parse(nonDefaultTemplate))

	if useDefaultVals == true {
		handler = cyoa.NewHandler(story)
	} else {
		handler = cyoa.NewHandler(story, cyoa.WithTemplate(tpl), cyoa.WithPath(nonDefaultPathFunc))
	}

	mux := http.NewServeMux()
	mux.Handle("/story/", handler)

	fmt.Printf("Starting server on port %d\n", *port)

	if useDefaultVals {
		log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), handler))
	} else {
		log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), mux))
	}
}
