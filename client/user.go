package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"strings"

	"github.com/Pallinder/go-randomdata"
	"github.com/jenkins-zh/jenkins-cli/util"
)

// UserClient for connect the user
type UserClient struct {
	JenkinsCore
}

// Token is the token of user
type Token struct {
	Status string
	Data   TokenData
}

// TokenData represents the token
type TokenData struct {
	TokenName  string
	TokenUUID  string
	TokenValue string
}

// Get returns a user's detail
func (q *UserClient) Get() (status *User, err error) {
	api := fmt.Sprintf("/user/%s/api/json", q.UserName)
	err = q.RequestWithData("GET", api, nil, nil, 200, &status)
	return
}

// EditDesc update the description of a user
func (q *UserClient) EditDesc(description string) (err error) {
	formData := url.Values{}
	formData.Add("description", description)
	payload := strings.NewReader(formData.Encode())
	_, err = q.RequestWithoutData("POST", fmt.Sprintf("/user/%s/submitDescription", q.UserName), map[string]string{util.ContentType: util.ApplicationForm}, payload, 200)
	return
}

// Delete will remove a user from Jenkins
func (q *UserClient) Delete(username string) (err error) {
	_, err = q.RequestWithoutData("POST", fmt.Sprintf("/securityRealm/user/%s/doDelete", username), map[string]string{util.ContentType: util.ApplicationForm}, nil, 200)
	return
}

func genSimpleUserAsPayload(username, password string) (payload io.Reader, user *UserForCreate) {
	user = &UserForCreate{
		User:      User{FullName: username},
		Username:  username,
		Password1: password,
		Password2: password,
		Email:     fmt.Sprintf("%s@%s.com", username, username),
	}

	userData, _ := json.Marshal(user)
	formData := url.Values{
		"json":      {string(userData)},
		"username":  {username},
		"password1": {password},
		"password2": {password},
		"fullname":  {username},
		"email":     {user.Email},
	}
	payload = strings.NewReader(formData.Encode())
	return
}

// Create will create a user in Jenkins
func (q *UserClient) Create(username, password string) (user *UserForCreate, err error) {
	var (
		payload io.Reader
		code int
	)

	if password == "" {
		password = util.GeneratePassword(8)
	}

	payload, user = genSimpleUserAsPayload(username, password)
	code, err = q.RequestWithoutData("POST", "/securityRealm/createAccountByAdmin",
		map[string]string{util.ContentType: util.ApplicationForm}, payload, 200)
	if code == 302 {
		err = nil
	}
	return
}

// CreateToken create a token in Jenkins
func (q *UserClient) CreateToken(newTokenName string) (status *Token, err error) {
	if newTokenName == "" {
		newTokenName = fmt.Sprintf("jcli-%s", randomdata.SillyName())
	}

	api := fmt.Sprintf("/user/%s/descriptorByName/jenkins.security.ApiTokenProperty/generateNewToken", q.UserName)

	formData := url.Values{}
	formData.Add("newTokenName", newTokenName)
	payload := strings.NewReader(formData.Encode())

	err = q.RequestWithData("POST", api,
		map[string]string{util.ContentType: util.ApplicationForm}, payload, 200, &status)
	return
}

// User for Jenkins
type User struct {
	AbsoluteURL string `json:"absoluteUrl"`
	Description string
	FullName    string `json:"fullname"`
	ID          string
}

// UserForCreate is the data for creatig a user
type UserForCreate struct {
	User      `json:",inline"`
	Username  string `json:"username"`
	Password1 string `json:"password1"`
	Password2 string `json:"password2"`
	Email     string `json:"email"`
}
