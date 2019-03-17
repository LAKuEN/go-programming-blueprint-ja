package main

import (
	"errors"
	"io/ioutil"
	"path/filepath"
)

// ErrNoAvatarURL はAvatarインスタンスがアバターのURLを返せないときに起こるエラー
var ErrNoAvatarURL = errors.New("chat: Cannot retrieve avatar image URL")

// Avatar はユーザーのプロフィール画像を表す
type Avatar interface {
	GetAvatarURL(ChatUser) (string, error)
}

type TryAvatars []Avatar

func (a TryAvatars) GetAvatarURL(u ChatUser) (string, error) {
	for _, avatar := range a {
		if url, err := avatar.GetAvatarURL(u); err == nil {
			return url, nil
		}
	}

	return "", ErrNoAvatarURL
}

type AuthAvatar struct{}

// UseAuthAvatar はAuthAvatarを使う箇所をわかり易くするための変数
var UseAuthAvatar AuthAvatar

func (AuthAvatar) GetAvatarURL(u ChatUser) (string, error) {
	url := u.AvatarURL()
	if url != "" {
		return url, nil
	}
	return "", ErrNoAvatarURL
}

// GravatarAvatar はGravatarからアバターを取得するための構造体
type GravatarAvatar struct{}

// UseGravatar はGravatarAvatarを使う箇所をわかり易くするための変数
var UseGravatar GravatarAvatar

func (GravatarAvatar) GetAvatarURL(u ChatUser) (string, error) {
	return "//www.gravatar.com/avatar/" + u.UniqueID(), nil
}

// FileSystemAvatar はローカルファイルとして保存しているアバターを取得するための構造体
type FileSystemAvatar struct{}

// UseFileSystemAvatar はFileSystemAvatarを使う場所をわかり易くするための変数
var UseFileSystemAvatar FileSystemAvatar

func (FileSystemAvatar) GetAvatarURL(u ChatUser) (string, error) {
	if files, err := ioutil.ReadDir("avatars"); err == nil {
		for _, file := range files {
			if file.IsDir() {
				continue
			}
			if match, _ := filepath.Match(u.UniqueID()+"*", file.Name()); match {
				return "/avatars/" + file.Name(), nil
			}
		}
	}

	return "", ErrNoAvatarURL
}
