package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"sync"
)

type templateHandler struct {
	once     sync.Once
	filename string
	tmpl     *template.Template
}

func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.once.Do(func() {
		tmplPath := filepath.Join("templates", t.filename)
		// https://golang.org/pkg/text/template/#Must
		t.tmpl = template.Must(template.ParseFiles(tmplPath))
	})

	if err := t.tmpl.Execute(w, r); err != nil {
		log.Fatal(err)
	}
}

func main() {
	addr := flag.String("addr", ":8080", "")
	logging := flag.Bool("logging", false, "")
	flag.Parse()

	http.Handle("/", &templateHandler{filename: "chat.html"})

	r := newRoom(*logging)
	http.Handle("/room", r)
	go r.run()

	fmt.Printf("Listening on port %v\n", *addr)
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatalf("Error occured at ListenAndServe: %#v", err)
	}
}
