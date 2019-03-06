package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"sync"

	"github.com/stretchr/gomniauth"
	"github.com/stretchr/gomniauth/providers/google"
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

	// 認証
	initGomniauth()

	http.Handle("/chat", MustAuth(&templateHandler{filename: "chat.html"}))
	http.Handle("/login", &templateHandler{filename: "login.html"})
	http.HandleFunc("/auth/", loginHandler)

	r := newRoom(*logging)
	http.Handle("/room", r)
	go r.run()

	fmt.Printf("Listening on port %v\n", *addr)
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatalf("Error occured at ListenAndServe: %#v", err)
	}
}

func initGomniauth() {
	// TODO 認証情報を渡す必要があるかも。以降の章進めるときに適宜直す
	type googleClientSecret struct {
		Web struct {
			ClientID     string `json:"client_id"`
			ClientSecret string `json:"client_secret"`
		} `json:"web"`
	}

	var loadGoogleClientSecret = func() (googleClientSecret, error) {
		jsonStr, err := ioutil.ReadFile("client_secret.apps.googleusercontent.com.json")
		if err != nil {
			return googleClientSecret{}, err
		}
		var clientSecret googleClientSecret
		err = json.Unmarshal(jsonStr, &clientSecret)
		if err != nil {
			return googleClientSecret{}, err
		}

		return clientSecret, nil
	}

	clientSecret, err := loadGoogleClientSecret()
	key, err := ioutil.ReadFile("security_key")
	if err != nil {
		log.Fatal(err)
	}

	gomniauth.SetSecurityKey(string(key))
	gomniauth.WithProviders(
		google.New(clientSecret.Web.ClientID, clientSecret.Web.ClientSecret, "http://localhost:8080/auth/callback/google"),
	)
}
