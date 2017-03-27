package models

import (
	"fmt"
	"golang.org/x/oauth2"
	"io/ioutil"
	"os"
	"path/filepath"
)

var (
	tokenDirectory = filepath.Join(os.Getenv("HOME"), ".kubecos")
	tokenFile      = "token"
)

func writeFile(str string) error {
	fmt.Println("TOKEN", str)
	err := os.MkdirAll(tokenDirectory, os.ModePerm)
	if err != nil {
		fmt.Println(err)
		return err
	}

	tokenPath := fmt.Sprintf("%s/%s", tokenDirectory, tokenFile)
	bts := []byte(str)
	err = ioutil.WriteFile(tokenPath, bts, 0644)
	fmt.Println(err)
	return err
}

//SaveAccessToken get access token from authorization code and saves locally
func SaveAccessToken(state, code, expectedState string) error {
	if state != expectedState {
		//fmt.Printf("invalid oauth state, expected '%s', got '%s'\n", oauthStateString, state)
		return nil
	}

	token, err := googleOauthConfig.Exchange(oauth2.NoContext, code)
	if err != nil {
		//fmt.Printf("Code exchange failed with '%s'\n", err)
		return nil
	}
	if !token.Valid() {
		return nil
	}

	err = writeFile(token.AccessToken)
	return err
}
