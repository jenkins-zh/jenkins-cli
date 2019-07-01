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

func (q *UserClient) Create(newTokenName string) (status *Token, err error) {
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

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
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
