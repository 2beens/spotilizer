package services

import (
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	m "github.com/2beens/spotilizer/models"
)

const requestTimeoutSeconds = 30

type httpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type requestCreator interface {
	NewRequest(method, url string, body io.Reader) (*http.Request, error)
}

type newRequestCreator struct{}

func (nrc newRequestCreator) NewRequest(method, url string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, url, body)
	return req, err
}

type requestClient struct {
	httpClient
	requestCreator
}

var reqClient = requestClient{
	// clients should be reused instead of created as needed
	// https://golang.org/pkg/net/http/#Client
	httpClient:     &http.Client{},
	requestCreator: newRequestCreator{},
}

func getFromSpotify(apiURL string, path string, authOptions *m.SpotifyAuthOptions) (body []byte, err error) {
	req, err := reqClient.requestCreator.NewRequest("GET", apiURL+path, nil)
	if err != nil {
		log.Printf(" >>> error getting spotify response. details: %v\n", err)
		return nil, err
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+authOptions.AccessToken)

	// complicating this function on purpose to demonstrante the usage of channels and goroutines
	// through implementing a request timeout mechanism
	timeoutChan := time.After(time.Duration(requestTimeoutSeconds) * time.Second)
	respChannel := make(chan []byte)
	errChannel := make(chan error)

	go func() {
		resp, reqErr := reqClient.httpClient.Do(req)
		if reqErr != nil {
			errChannel <- reqErr
			return
		}
		body, reqErr = ioutil.ReadAll(resp.Body)
		if reqErr != nil {
			errChannel <- reqErr
			return
		}
		respChannel <- body
	}()

	select {
	case <-timeoutChan:
		return nil, errors.New("timeout occured")
	case body = <-respChannel:
		return body, nil
	case err = <-errChannel:
		return nil, err
	}
}

func getAPIError(body []byte) (spErr m.SpAPIError, isError bool) {
	err := json.Unmarshal(body, &spErr)
	if err == nil && len(spErr.Error.Message) > 0 && spErr.Error.Status > 0 {
		return spErr, true
	}
	return m.SpAPIError{}, false
}
