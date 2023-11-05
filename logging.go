package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"time"
)

type LoggingInterface interface {
	doRequest(req *http.Request) (*http.Response, error)
}

type LoggingService struct {
	next LoggingInterface
}

// log the interface request handler
func (s *LoggingService) doRequest(request *http.Request) (response *http.Response, err error) {
	var responseBody []byte
	defer func(start time.Time) {
		fmt.Printf("request=%v, err=%v, took=%v, response=%s\n", request.URL, err, time.Since(start), responseBody)
	}(time.Now())
	resp, err := s.next.doRequest(request)
	responseBody, err = io.ReadAll(resp.Body)
	resp.Body = io.NopCloser(bytes.NewBuffer(responseBody))
	return resp, err
}

func NewLoggingService(next LoggingInterface) LoggingInterface {
	return &LoggingService{next: next}
}
