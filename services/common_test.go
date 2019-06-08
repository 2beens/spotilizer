package services

import (
	"bytes"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type httpClientMock struct{}

var testURL = "test-url"
var testPathOK = "/test/path/ok"
var testPathErr = "/test/path/err"
var testPathTimeout = "/test/path/timeout"
var errorPathMessage = "dummy-response-error"
var accessToken = "test-accessToken"

func (c httpClientMock) Do(req *http.Request) (*http.Response, error) {
	log.Println(" > http client mock, Do(req) path: " + req.URL.Path)
	switch req.URL.Path {
	case testURL + testPathTimeout:
		time.Sleep(time.Duration(reqClient.requestTimeoutSeconds+5) * time.Second)
		return &http.Response{
			Body:       ioutil.NopCloser(bytes.NewBufferString("OK")),
			StatusCode: 200,
		}, nil
	case testURL + testPathOK:
		return &http.Response{
			Body:       ioutil.NopCloser(bytes.NewBufferString("OK")),
			StatusCode: 200,
		}, nil
	case testURL + testPathErr:
		return nil, errors.New(errorPathMessage)
	default:
		return nil, errors.New("this is not intended to be reached")
	}
}

func TestGetFromSpotify(t *testing.T) {
	log.Println(" > TestGetFromSpotify: starting ...")
	reqClient = requestClient{httpClient: &httpClientMock{}, requestTimeoutSeconds: 1}

	body, err := getFromSpotify(testURL, testPathOK, accessToken)
	assert.NoError(t, err)
	if assert.Equal(t, "OK", string(body), "response body not correct") {
		log.Println(" > no error response OK")
	}

	body, err = getFromSpotify(testURL, testPathErr, accessToken)
	assert.EqualError(t, err, errorPathMessage)
	if assert.Equal(t, []byte(nil), body, "response body not correct") {
		log.Println(" > error response OK")
	}

	body, err = getFromSpotify(testURL, testPathTimeout, accessToken)
	assert.EqualError(t, err, "timeout occured")
	if assert.Equal(t, []byte(nil), body, "response body not correct") {
		log.Println(" > timeout, error response OK")
	}

	log.Println(" > TestGetFromSpotify: tests finished!")
}