package client

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/Pallinder/go-randomdata"
	"github.com/jenkins-zh/jenkins-cli/util"
)

type UserClient struct {
	JenkinsCore
}

type Token struct {
	Status string
	Data   TokenData
}

type TokenData struct {
	TokenName  string
	TokenUuid  string
	TokenValue string
}

func (q *UserClient) Get() (status *User, err error) {
	api := fmt.Sprintf("%s/user/%s/api/json", q.URL, q.UserName)
	var (
		req      *http.Request
		response *http.Response
	)

	req, err = http.NewRequest("GET", api, nil)
	if err == nil {
		q.AuthHandle(req)
	} else {
		return
	}

	client := q.GetClient()
	if response, err = client.Do(req); err == nil {
		code := response.StatusCode
		var data []byte
		data, err = ioutil.ReadAll(response.Body)
		if code == 200 {
			if err == nil {
				status = &User{}
				err = json.Unmarshal(data, status)
			}
		} else {
			log.Fatal(string(data))
		}
	} else {
		log.Fatal(err)
	}
	return
}

func (q *UserClient) EditDesc(description string) (err error) {
	api := fmt.Sprintf("%s/user/%s/submitDescription", q.URL, q.UserName)
	var (
		req      *http.Request
		response *http.Response
	)

	formData := url.Values{}
	formData.Add("description", description)
	payload := strings.NewReader(formData.Encode())

	req, err = http.NewRequest("POST", api, payload)
	if err == nil {
		q.AuthHandle(req)
	} else {
		return
	}

	req.Header.Set(util.CONTENT_TYPE, util.APP_FORM)
	if err = q.CrumbHandle(req); err != nil {
		log.Fatal(err)
	}

	client := q.GetClient()
	if response, err = client.Do(req); err == nil {
		code := response.StatusCode
		var data []byte
		data, err = ioutil.ReadAll(response.Body)
		if code == 200 {
		} else {
			log.Fatal(string(data))
		}
	} else {
		log.Fatal(err)
	}
	return
}

func (q *UserClient) Delete(username string) (err error) {
	api := fmt.Sprintf("%s/securityRealm/user/%s/doDelete", q.URL, username)
	var (
		req      *http.Request
		response *http.Response
	)

	req, err = http.NewRequest("POST", api, nil)
	if err == nil {
		q.AuthHandle(req)
	} else {
		return
	}

	req.Header.Set(util.CONTENT_TYPE, util.APP_FORM)
	client := q.GetClient()
	if response, err = client.Do(req); err == nil {
		code := response.StatusCode
		var data []byte
		data, err = ioutil.ReadAll(response.Body)
		if code != 200 && code != 302 {
			log.Fatal(string(data))
		}
	} else {
		log.Fatal(err)
	}
	return
}

func (q *UserClient) Create(username string) (user *UserForCreate, err error) {
	api := fmt.Sprintf("%s/securityRealm/createAccountByAdmin", q.URL)
	var (
		req      *http.Request
		response *http.Response
	)

	passwd := util.GeneratePassword(8)

	user = &UserForCreate{
		User:      User{FullName: username},
		Username:  username,
		Password1: passwd,
		Password2: passwd,
		Email:     fmt.Sprintf("%s@%s.com", username, username),
	}

	userData, _ := json.Marshal(user)
	formData := url.Values{
		"json":      {string(userData)},
		"username":  {username},
		"password1": {passwd},
		"password2": {passwd},
		"fullname":  {username},
		"email":     {user.Email},
	}
	payload := strings.NewReader(formData.Encode())

	req, err = http.NewRequest("POST", api, payload)
	if err == nil {
		q.AuthHandle(req)
	} else {
		return
	}

	req.Header.Set(util.CONTENT_TYPE, util.APP_FORM)
	client := q.GetClient()
	if response, err = client.Do(req); err == nil {
		code := response.StatusCode
		var data []byte
		data, err = ioutil.ReadAll(response.Body)
		if code != 200 && code != 302 {
			log.Fatal(string(data))
		}
	} else {
		log.Fatal(err)
	}
	return
}

func (q *UserClient) CreateToken(newTokenName string) (status *Token, err error) {
	if newTokenName == "" {
		newTokenName = fmt.Sprintf("jcli-%s", randomdata.SillyName())
	}

	api := fmt.Sprintf("%s/user/%s/descriptorByName/jenkins.security.ApiTokenProperty/generateNewToken", q.URL, q.UserName)
	var (
		req      *http.Request
		response *http.Response
	)

	formData := url.Values{}
	formData.Add("newTokenName", newTokenName)
	payload := strings.NewReader(formData.Encode())

	req, err = http.NewRequest("POST", api, payload)
	if err == nil {
		q.AuthHandle(req)
	} else {
		return
	}

	req.Header.Set(util.CONTENT_TYPE, util.APP_FORM)
	if err = q.CrumbHandle(req); err != nil {
		log.Fatal(err)
	}

	client := q.GetClient()
	if response, err = client.Do(req); err == nil {
		code := response.StatusCode
		var data []byte
		data, err = ioutil.ReadAll(response.Body)
		if code == 200 {
			if err == nil {
				status = &Token{}
				err = json.Unmarshal(data, status)
			}
		} else {
			log.Fatal(string(data))
		}
	} else {
		log.Fatal(err)
	}
	return
}

type User struct {
	AbsoluteURL string `json:"absoluteUrl"`
	Description string
	FullName    string `json:"fullname"`
	ID          string
}

type UserForCreate struct {
	User      `json:inline`
	Username  string `json:"username"`
	Password1 string `json:"password1"`
	Password2 string `json:"password2"`
	Email     string `json:"email"`
}
