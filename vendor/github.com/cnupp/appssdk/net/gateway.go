package net

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/cnupp/appssdk/config"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
)

type Gateway struct {
	config config.Reader
}

func NewCloudControllerGateway(reader config.Reader) Gateway {
	return Gateway{config: reader}
}

const (
	contentType string = "application/json"
	userAgent   string = "builder"
)

type Request struct {
	HttpReq      *http.Request
	SeekableBody io.ReadSeeker
}

func (g *Gateway) Get(path string, value interface{}) (getErr error) {
	res, err := g.Request("GET", path, nil)
	if err != nil {
		getErr = err
		return
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		getErr = err
		return
	}

	getErr = json.Unmarshal(body, value)
	return
}

func (g *Gateway) Delete(path string, value interface{}) (getErr error) {
	res, err := g.Request("DELETE", path, nil)
	if err != nil {
		getErr = err
		return
	}

	_, err = ioutil.ReadAll(res.Body)
	if err != nil {
		getErr = err
		return
	}

	return
}

func (g *Gateway) PUT(path string, value interface{}) (apiErr error) {
	data, err := json.Marshal(value)
	if err != nil {
		apiErr = fmt.Errorf("Can not serilize the data")
		return
	}
	_, apiErr = g.Request("PUT", path, data)
	return
}

func (g *Gateway) Request(method string, path string, body []byte) (*http.Response, error) {
	return g.request(method, path, body, contentType)
}

func checkForErrors(res *http.Response, body string) error {

	// If response is not an error, return nil.
	if res.StatusCode > 199 && res.StatusCode < 400 {
		return nil
	}

	// Read the response body if none was provided.
	if body == "" {
		defer res.Body.Close()
		resBody, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return err
		}
		body = string(resBody)
	}

	// Unmarshal the response as JSON, or return the status and body.
	bodyMap := make(map[string]interface{})
	if err := json.Unmarshal([]byte(body), &bodyMap); err != nil {
		return fmt.Errorf("\n%s\n%s\n", res.Status, body)
	}

	errorMessage := fmt.Sprintf("\n%s\n", res.Status)
	for key, value := range bodyMap {
		switch v := value.(type) {
		case string:
			errorMessage += fmt.Sprintf("%s: %s\n", key, v)
		case []interface{}:
			for _, subValue := range v {
				switch sv := subValue.(type) {
				case string:
					errorMessage += fmt.Sprintf("%s: %s\n", key, sv)
				case map[string]interface{}:
					errorMessage += fmt.Sprintf("%s\n", sv["detail"])
				default:
					fmt.Printf("Unexpected type in %s error message array. Contents: %v",
						reflect.TypeOf(value), sv)
				}
			}
		default:
			fmt.Printf("Cannot handle key %s in error message, type %s. Contents: %v",
				key, reflect.TypeOf(value), bodyMap[key])
		}
	}

	return errors.New(errorMessage)
}

func (g *Gateway) request(method string, path string, body []byte, contentType string) (*http.Response, error) {
	u, err := url.Parse(path)
	if err != nil {
		return nil, err
	}

	var uri string
	if "" == u.Scheme {
		uri = g.config.ApiEndpoint() + path
	} else {
		uri = path
	}
	req, err := http.NewRequest(method, uri, bytes.NewBuffer(body))

	if err != nil {
		return nil, err
	}
	req.Close = true
	req.Header.Add("Content-Type", contentType)
	req.Header.Add("Accept", contentType)
	req.Header.Add("User-Agent", userAgent)
	req.Header.Add("Authorization", g.config.Auth())

	client := http.Client{}

	res, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	if err = checkForErrors(res, ""); err != nil {
		return nil, err
	}

	return res, nil
}
