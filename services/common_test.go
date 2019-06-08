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

func TestGetAPIError(t *testing.T) {
	log.Println(" > TestGetAPIError: starting ...")

	apiErrJSON := []byte(`{"error": {"status": 401, "message": "Test API Error"}}`)
	apiErr, isErr := getAPIError(apiErrJSON)
	assert.True(t, isErr, "API Response not received")
	if assert.NotNil(t, apiErr, "API Error should not be nil") {
		log.Println(" > API error response not nil, OK")
	}
	assert.Equal(t, "Test API Error", apiErr.Error.Message, "API Error message not correct")
	assert.Equal(t, 401, apiErr.Error.Status, "API Error status not correct")

	apiErrJSONWrong1 := []byte(`{"status": 401, "message": "Test API Error"}`)
	apiErr, isErr = getAPIError(apiErrJSONWrong1)
	if assert.False(t, isErr, "API Response 1, err should not be false") {
		log.Println(" > API error 1, wrong JSON, isErr false, OK")
	}

	apiErrJSONWrong2 := []byte(`some text which is not JSON at all`)
	apiErr, isErr = getAPIError(apiErrJSONWrong2)
	if assert.False(t, isErr, "API Response 2, err should not be false") {
		log.Println(" > API error 2, wrong JSON, isErr false, OK")
	}

	log.Println(" > TestGetAPIError: tests finished!")
}
