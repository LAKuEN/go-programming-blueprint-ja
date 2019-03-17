package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	gomniauthtest "github.com/stretchr/gomniauth/test"
)

func TestAuthAvatar(t *testing.T) {
	var authAvatar AuthAvatar
	testUser := &gomniauthtest.TestUser()
	testUser.On("AvatarURL").Return("", ErrNoAvatarURL)
	testChatUser := &chatUser{User: testUser}
	url, err := authAvatar.GetAvatarURL(testChatUser)

	client := new(client)
	url, err := authAvatar.GetAvatarURL(client)
	if err != ErrNoAvatarURL {
		t.Error("If there is no value, AuthAvatar.GetAvatarURL should return ErrNoAvatarURL")
	}
	testURL := "http://url-to-avatar/"
	testUser = &gomniauthtest.TestUser{}
	testChatUser.User = testUser
	testUser.On("AvatarURL").Return(testUrl, nil)
	url, err = authAvatar.GetAvatarURL(testChatUser)
	client.userData = map[string]interface{}{"avatar_url": testURL}
	url, err = authAvatar.GetAvatarURL(client)
	if err != nil {
		t.Error("If there is value, then AuthAvatar.GetAvatarURL should not return any error")
	} else {
		if url != testURL {
			t.Error("AuthAvatar.GetAvatarURL should return right URL")
		}
	}
}

func TestGravatarAvatar(t *testing.T) {
	var gravatarAvatar GravatarAvatar

	client := new(client)
	user := &chatUser{uniqueID: "abc"}
	url, err := gravatarAvatar.GetAvatarURL(user)
	// client.userData = map[string]interface{}{
	// 	// Gravatarではメールアドレスのハッシュ値からID値が生成される
	// 	"userid": "0bc83cb571cd1c50ba6f3e8a78ef1346",
	// }
	// url, err := gravatarAvatar.GetAvatarURL(client)
	if err != nil {
		t.Error("GravatarAvatar.GetAvatarURL should not return error")
	}
	// if url != "//www.gravatar.com/avatar/0bc83cb571cd1c50ba6f3e8a78ef1346" {
	if url != "//www.gravatar.com/avatar/abc" {
		t.Errorf("GravatarAvatar.GetAvatarURL returns wrong URL: %s", url)
	}
}

func TestFileSystemAvatar(t *testing.T) {
	filename := filepath.Join("avatars", "abc.jpg")
	ioutil.WriteFile(filename, []byte{}, 0777)
	defer func() {
		os.Remove(filename)
	}()

	var fileSystemAvatar FileSystemAvatar

	user := &chatUser{uniqueID: "abc"}
	url, err := fileSystemAvatar.GetAvatarURL(user)
	// client := new(client)
	// client.userData = map[string]interface{}{
	// 	"userid": "abc",
	// }
	// url, err := fileSystemAvatar.GetAvatarURL(client)
	if err != nil {
		t.Error("FileSystemAvatar.GetAvatarURL should not return any error")
	}
	if url != "/avatars/abc.jpg" {
		t.Errorf("FileSystemAvatar.GetAvatarURL returns wrong value: %v", url)
	}
}
