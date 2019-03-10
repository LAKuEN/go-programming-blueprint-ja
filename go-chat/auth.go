package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/stretchr/gomniauth"
	"github.com/stretchr/objx"
)

type authHandler struct {
	next http.Handler
}

func (h *authHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if _, err := r.Cookie("auth"); err == http.ErrNoCookie {
		w.Header().Set("Location", "/login")
		w.WriteHeader(http.StatusTemporaryRedirect)
	} else if err != nil {
		panic(err.Error())
	} else {
		h.next.ServeHTTP(w, r)
	}
}

// MustAuth はhttp.Handlerインタフェースに適合したハンドラをラップしたハンドラを返却する
func MustAuth(handler http.Handler) http.Handler {
	return &authHandler{next: handler}
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	segs := strings.Split(r.URL.Path, "/")
	if len(segs) < 4 {
		panic(fmt.Errorf("Specified invalid URI: %v", r.URL.Path))
	}
	action := segs[2]
	providerName := segs[3]

	switch action {
	case "login":
		provider, err := gomniauth.Provider(providerName)
		if err != nil {
			log.Fatalf("Failed to get authentication provider: %v-%v", provider, err)
		}
		loginURL, err := provider.GetBeginAuthURL(nil, nil)
		if err != nil {
			log.Fatalf("Error occured when called GetBeginAuthURL: %v-%v", provider, err)
		}
		w.Header().Set("Location", loginURL)
		w.WriteHeader(http.StatusTemporaryRedirect)
	case "callback":
		provider, err := gomniauth.Provider(providerName)
		if err != nil {
			log.Fatalf("Failed to get authentication provider: %v-%v", provider, err)
		}
		authCode := objx.MustFromURLQuery(r.URL.RawQuery)
		// 認可コード -> アクセストークン
		creds, err := provider.CompleteAuth(authCode)
		if err != nil {
			log.Fatalf("Couldn't finish authentication: %v", err)
		}
		// アクセストークン -> ユーザ情報
		user, err := provider.GetUser(creds)
		if err != nil {
			log.Fatalf("Couldn't get user information: %v", err)
		}
		// ユーザ情報 -> Cookie
		authCookieValue := objx.New(map[string]interface{}{
			"name":       user.Name(),
			"avatar_url": user.AvatarURL(), // アバター画像の格納先はサービスに依って異なるが、gomniauthが差異を吸収してくれる
		}).MustBase64()
		http.SetCookie(w, &http.Cookie{
			Name:  "auth",
			Value: authCookieValue,
			Path:  "/",
		})

		w.Header()["Location"] = []string{"/chat"}
		w.WriteHeader(http.StatusTemporaryRedirect)
	default:
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "Not defined action: %v", action)
	}
}
