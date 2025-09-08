// Package restclient provides a REST client to communicate with the Liana Technologies REST API.
package restclient

import (
	"fmt"
	"net/http"
	"encoding/json"
	"encoding/hex"
	"strings"
	"crypto/md5"
	"crypto/sha256"
	"crypto/hmac"
	"time"
	"io/ioutil"
	"errors"
)

type restclient struct {

	apiUser	int
	apiSecret string
	apiUrl string
	apiVersion int
	apiRealm string

	httpclient *http.Client

	endpoint string
	bodystr string
	hashstr string
	timestr string
}

func NewRestClient(apiUser int, apiSecret string, apiUrl string, apiVersion int, apiRealm string) (restclient, error) {

	// Create http client
	httpclient := &http.Client{
		Timeout: 60 * time.Second,
	}

	return restclient {apiUser, apiSecret, apiUrl, apiVersion, apiRealm, httpclient, "", "", "" ,""}, nil
}

// Call is used to perform a call to the RESTful API.
func (this restclient) Call(path string, params interface{}, inputMethod ...string) (interface{}, error) {
	bodyjson, err := json.Marshal(params)

	if err != nil {
		return nil, err
	}

	this.bodystr = string(bodyjson)

	method := "POST"
	if len(inputMethod) > 0 {
		method = inputMethod[0]
	}
	if method == "GET" {
		this.bodystr = "";
	}

	hash := md5.New()
	hash.Write([]byte(string(this.bodystr)))
	this.hashstr = hex.EncodeToString(hash.Sum(nil))
	this.endpoint = fmt.Sprintf("/api/v%d/%s", this.apiVersion, path)
	this.timestr = time.Now().Format(time.RFC3339)

	req, err := this.createRequest(method)

	if err != nil {
		return nil, err
	}

	resp, err := this.httpclient.Do(req)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	return this.handleResponse(resp)
}

func (this restclient) sign(sign_string string) (interface{}, error) {
	sign := hmac.New(sha256.New, []byte(this.apiSecret))
	sign.Write([]byte(sign_string))
	return hex.EncodeToString(sign.Sum(nil)), nil
}

// Form the http.Request to fetch the response data
func (this restclient) createRequest(method string) (*http.Request, error) {
	req, err := http.NewRequest(
		method,
		this.apiUrl + this.endpoint,
		strings.NewReader(this.bodystr),
	)

	if err != nil {
		return nil, err
	}

	signature, err := this.sign(fmt.Sprintf("%s\n%s\n%s\n%s\n%s\n%s",
		method,
		this.hashstr,
		"application/json",
		this.timestr,
		this.bodystr,
		this.endpoint,
	))

	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Content-MD5", this.hashstr)
	req.Header.Add("Date", this.timestr)

	req.Header.Add(
		"Authorization",
		fmt.Sprintf("%s %d:%s", this.apiRealm, this.apiUser, signature),
	)

	return req, nil
}

// Preprocesses an API response
// Response json is decoded and checked for the success value.
func (this restclient) handleResponse(resp *http.Response) (interface{}, error) {
	var response map[string]interface{}

	bodystr, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	json.Unmarshal(bodystr, &response)

	if this.apiVersion == 1 || this.apiVersion == 2 {
		var succeed, ok = response["succeed"]
		if ! ok {
			return nil, errors.New("unexpected response from API: " + string(bodystr))
		}
		if ! succeed.(bool) {
			return nil, errors.New(this.endpoint + ": " + response["message"].(string))
		}
		return response["result"], nil
	} else if this.apiVersion == 3 {
		var _, ok = response["items"]
		if ! ok {
			return nil, errors.New("unexpected response from API: " + string(bodystr))
		}
		return response["items"], nil
	}

	return nil, errors.New(fmt.Sprintf("unexpected api version %+v", this.apiVersion))
}

