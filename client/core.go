package client

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type CoreClient struct {
	JenkinsCore
}

// Search find a set of jobs by name
func (q *CoreClient) Restart() (err error) {
	api := fmt.Sprintf("%s/safeRestart", q.URL)
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

	client := q.GetClient()
	if response, err = client.Do(req); err == nil {
		code := response.StatusCode
		var data []byte
		data, err = ioutil.ReadAll(response.Body)
		if code == 503 { // Jenkins could be behind of a proxy
			fmt.Println("Please wait while Jenkins is restarting")
		} else if code != 200 || err != nil {
			log.Fatalf("Error code: %d, response: %s, errror: %v", code, string(data), err)
		} else {
			fmt.Println("restart successfully")
		}
	} else {
		log.Fatal(err)
	}
	return
}
