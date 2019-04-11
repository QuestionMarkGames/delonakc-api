package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
)

type Github struct {
	BaseUrl string `yaml:"baseUrl"`
	TokenUrl string `yaml:"tokenUrl"`
	InfoUrl string 	`yaml:"infoUrl"`
	ClientId string `yaml:"clientId"`
	ClientSecret string `yaml:"clientSecret"`
	Login string `yaml:"login"`
	Token string
	Code string
	RedirectUrl string
	State string
}

type GithubToken struct {
	AccessToken string `json:"access_token""`
	TokenType string `json:"token_type"`
	Scope string `json:"scope"`
}

type GithubUser struct {
	Name string `json:"name"`
	AvatarUrl string `json:"avatar_url"`
	Id int	`json:"id"`
	Home string `json:"html_url"`
}

func (g *Github) ExchangeToken() (*GithubToken, error)  {

	log.SetPrefix("<GithubAuth>")
	defer log.SetPrefix("")

	params := &url.Values{}
	params.Add("client_id", g.ClientId)
	params.Add("client_secret", g.ClientSecret)
	params.Add("code", g.Code)
	params.Add("redirect_url", g.RedirectUrl)
	params.Add("state", g.State)

	body := strings.NewReader(params.Encode())
	req, err := http.NewRequest("POST", g.TokenUrl, body)

	if err != nil {
		log.Println("error with new request: ", err)
		return nil, err
	}

	req.Header.Set("Accept", "application/json")

	client := &http.Client{}

	res, err := client.Do(req)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	resBody, _ := ioutil.ReadAll(res.Body)

	var token *GithubToken

	err = json.Unmarshal(resBody, &token)

	if err != nil {
		return nil, err
	}

	g.Token = token.AccessToken

	return token, nil
}

func (g *Github) GetUserInfo() (*GithubUser, error) {
	req, err := http.NewRequest("GET", g.InfoUrl, nil)

	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", fmt.Sprintf("token %s", g.Token))

	client := &http.Client{}

	res, err := client.Do(req)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	var user *GithubUser

	resBody, _ := ioutil.ReadAll(res.Body)

	//fmt.Printf("%s\n", resBody)

	err = json.Unmarshal(resBody, &user)

	if err != nil {
		return nil, err
	}

	return user, nil
}