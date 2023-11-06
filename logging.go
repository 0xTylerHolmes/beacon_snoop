package main

import (
	"bytes"
	"fmt"
	"github.com/rs/zerolog"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

type LoggingInterface interface {
	doRequest(req *http.Request) (*http.Response, error)
}

type Header struct {
	Headers map[string][]string
}

// Impl for Stringer Interface
func (h *Header) String() string {
	s := ""
	for header, value := range h.Headers {
		s += fmt.Sprintf("%s: %s,", header, value)
	}
	return strings.TrimSuffix(s, ",")
}

// requestDumpBody dumps the request body, headers and replenishes the body reader
func dumpRequest(request *http.Request) ([]byte, *Header, error) {
	headers := Header{Headers: make(map[string][]string)}
	for key, value := range request.Header {
		headers.Headers[key] = value
	}
	if request.Method == "POST" {
		requestBody, err := io.ReadAll(request.Body)
		if err != nil {
			return nil, &headers, err
		}
		request.Body = io.NopCloser(bytes.NewBuffer(requestBody))
		return requestBody, &headers, nil
	}
	return nil, &headers, nil

}

func dumpResponse(response *http.Response) ([]byte, *Header, error) {
	responseBody, err := io.ReadAll(response.Body)
	headers := Header{Headers: make(map[string][]string)}
	if err != nil {
		return nil, nil, err
	}
	for key, value := range response.Header {
		headers.Headers[key] = value
	}
	response.Body = io.NopCloser(bytes.NewBuffer(responseBody))
	return responseBody, &headers, nil
}

type LoggingService struct {
	next LoggingInterface
	// if not nil logs to file. More verbose than the stdout logger
	fileLogger *zerolog.Logger
	// whether we should log the headers (for debugging)
	logHeaders bool
}

// log the interface request handler
func (s *LoggingService) doRequest(request *http.Request) (response *http.Response, err error) {
	var responseBody []byte
	var responseHeaders *Header
	var requestHeaders *Header
	requestBody, requestHeaders, err := dumpRequest(request)
	startTime := time.Now()
	resp, requestErr := s.next.doRequest(request)
	duration := time.Since(startTime)
	if requestErr == nil {
		responseBody, responseHeaders, err = dumpResponse(resp)
	}
	fmt.Printf("request=%v, request_body=%s, err=%v, took=%v, response=%s\n", request.URL, requestBody, err, duration, responseBody)
	if s.fileLogger != nil {
		if s.logHeaders {
			s.fileLogger.Info().Str("request_url", request.URL.String()).Str("request_type", request.Method).Stringer("request_headers", requestHeaders).Str("request_body", string(requestBody)).Err(requestErr).Str("took", duration.String()).Stringer("response_headers", responseHeaders).Str("response", string(responseBody)).Msg("")
		} else {
			s.fileLogger.Info().Str("request_url", request.URL.String()).Str("request_type", request.Method).Str("request_body", string(requestBody)).Err(requestErr).Str("took", duration.String()).Str("response", string(responseBody)).Msg("")
		}
	}
	return resp, requestErr
}

// NewLoggingService returns a loggingInterface. If logFileWriter is nil then we don't log to file
func NewLoggingService(next LoggingInterface, logFileWriter *os.File, logHeaders bool) LoggingInterface {
	service := &LoggingService{
		next:       next,
		fileLogger: nil,
		logHeaders: logHeaders,
	}
	if logFileWriter != nil {
		fileLogger := zerolog.New(logFileWriter).With().Timestamp().Logger().Output(logFileWriter).Level(zerolog.InfoLevel)
		service.fileLogger = &fileLogger
	}
	return service
}
