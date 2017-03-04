package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/webapps/ridemerge"
)

// serveTemplate parses template files and serves resulting html.
func serveTemplate(w http.ResponseWriter, session ridemerge.Session, filenames ...string) {
	var files []string
	for _, file := range filenames {
		files = append(files, fmt.Sprintf("templates/%s.html", file))
	}

	t := template.Must(template.ParseFiles(files...))
	err := t.ExecuteTemplate(w, "layout", session)
	if err != nil {
		log.Print("template executing error: ", err)
	}
}
