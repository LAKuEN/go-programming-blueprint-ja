package main

import (
	"io"
	"io/ioutil"
	"net/http"
	"path/filepath"
)

func uploaderHandler(w http.ResponseWriter, req *http.Request) {
	userID := req.FormValue("userid")
	file, header, err := req.FormFile("avatarFile")
	if err != nil {
		io.WriteString(w, err.Error())
		return
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		io.WriteString(w, err.Error())
		return
	}

	// アバター画像を保存
	filename := filepath.Join("avatars", userID+filepath.Ext(header.Filename))
	err = ioutil.WriteFile(filename, data, 0644)
	if err != nil {
		io.WriteString(w, err.Error())
		return
	}

	io.WriteString(w, "Success")
}
