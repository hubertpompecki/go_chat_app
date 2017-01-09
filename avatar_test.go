package main

import (
    "testing"
)

func TestAuthAvatar(t *testing.T) {
    var authAvatar AuthAvatar
    client := new(client)

    // no value
    url, err := authAvatar.GetAvatarURL(client)
    if err != ErrNoAvatarURL {
        t.Error("AuthAvatar.GetAvatarURL should return ErrNoAvatarURL when no value present")
    }

    // set value
    testUrl := "http://url-to-gravatar/"
    client.userData = map[string]interface{} {"avatar_url": testUrl}
    url, err = authAvatar.GetAvatarURL(client)
    if err != nil {
        t.Error("AuthAvatar.GetAvatarURL should return no error when value present")
    } else {
        if url != testUrl {
            t.Error("AuthAvatar.GetAvatarURL should return correct URL")
        }
    }
}