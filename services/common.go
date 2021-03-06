package services

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/2beens/spotilizer/models"
	log "github.com/sirupsen/logrus"
)

type httpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type requestClient struct {
	httpClient
	requestTimeoutSeconds int
}

var reqClient = requestClient{
	// clients should be reused instead of created as needed
	// https://golang.org/pkg/net/http/#Client
	httpClient:            &http.Client{},
	requestTimeoutSeconds: 30,
}

func getFromSpotify(apiURL string, path string, accessToken string) (body []byte, err error) {
	log.Tracef(" > getting from Spotify API [%s]: %s", apiURL, path)
	req, err := http.NewRequest("GET", apiURL+path, nil)
	if err != nil {
		log.Infof(" >>> error getting spotify response. details: %s", err.Error())
		return nil, err
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+accessToken)

	// complicating this function on purpose to demonstrate the usage of channels and goroutines
	// through implementing a request timeout mechanism
	timeoutChan := time.After(time.Duration(reqClient.requestTimeoutSeconds) * time.Second)
	respChannel := make(chan []byte)
	errChannel := make(chan error)

	go func() {
		resp, reqErr := reqClient.httpClient.Do(req)
		if reqErr != nil {
			errChannel <- reqErr
			return
		}
		defer resp.Body.Close()

		body, reqErr = ioutil.ReadAll(resp.Body)
		if reqErr != nil {
			errChannel <- reqErr
			return
		}
		respChannel <- body
	}()

	select {
	case <-timeoutChan:
		return nil, errors.New("timeout occurred")
	case body = <-respChannel:
		return body, nil
	case err = <-errChannel:
		return nil, err
	}
}

func getAPIError(body []byte) (spErr models.SpAPIError, isError bool) {
	err := json.Unmarshal(body, &spErr)
	if err == nil && len(spErr.Error.Message) > 0 && spErr.Error.Status > 0 {
		return spErr, true
	}
	return models.SpAPIError{}, false
}
