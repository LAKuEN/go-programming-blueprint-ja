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
	"github.com/stretchr/objx"
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
	data := map[string]interface{}{
		"Host": r.Host,
	}
	if authCookie, err := r.Cookie("auth"); err == nil {
		data["UserData"] = objx.MustFromBase64(authCookie.Value)
	}

	if err := t.tmpl.Execute(w, data); err != nil {
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
	http.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {
		// Cookieの削除
		// MaxAgeを-1にすると即座に削除されるが、ブラウザに依って削除されないケースがあるので
		// 空文字列をはめ込む
		http.SetCookie(w, &http.Cookie{
			Name:   "auth",
			Value:  "",
			Path:   "/",
			MaxAge: -1,
		})

		w.Header()["Location"] = []string{"/chat"}
		w.WriteHeader(http.StatusTemporaryRedirect)
	})
	http.HandleFunc("/auth/", loginHandler) // 末尾に/をつけると接頭辞として扱われる

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
