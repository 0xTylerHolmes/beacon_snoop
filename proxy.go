package beacon_snoop

import (
	"bytes"
	"context"
	"fmt"
	"github.com/pkg/errors"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
)

type SnooperConfig struct {
	remote     string
	listenAddr string
}

type Snooper struct {
	config       SnooperConfig
	target       *url.URL
	server       *http.Server
	client       *http.Client
	logger       LoggingInterface
	reverseProxy *httputil.ReverseProxy
}

// populate the headers of a new request with only keys we care about
func populateHeaders(source http.Request, dest *http.Request) {
	for k, v := range source.Header {
		if k == "Accept" || k == "Content-Type" || k == "User-Agent" || k == "Eth-Consensus-Version" {
			for _, value := range v {
				dest.Header.Add(k, value)
			}
		}
	}
}

// doRequest split here to allow us to benchmark via logginginterface
func (s *Snooper) doRequest(req *http.Request) (*http.Response, error) {
	return s.client.Do(req)
}

// RoundTrip implements the round trip used by the reverseProxy ServeHTTP()
func (s *Snooper) RoundTrip(request *http.Request) (*http.Response, error) {

	newRequest, err := http.NewRequest(request.Method, request.URL.String(), request.Body)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create proxy request")
	}
	populateHeaders(*request, newRequest)

	response, err := s.logger.doRequest(newRequest) // with logging
	//response, err := s.doRequest(newRequest)  // without logging
	if err != nil {
		// failed on request handling
		return nil, err
	}
	data, err := io.ReadAll(response.Body)
	response.Body = io.NopCloser(bytes.NewBuffer(data))
	return response, err
}

func (s *Snooper) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.reverseProxy.ServeHTTP(w, r)
}

func NewSnooper(config SnooperConfig) (*Snooper, error) {
	target, err := url.Parse(config.remote)
	if err != nil {
		return nil, err
	}
	snooper := &Snooper{
		config: config,
		target: target,
	}
	// proxy to use for serving requests
	snooper.reverseProxy = httputil.NewSingleHostReverseProxy(target)
	snooper.reverseProxy.Transport = snooper
	snooper.reverseProxy.Rewrite = nil
	snooper.logger = NewLoggingService(snooper)
	snooper.server = &http.Server{
		Addr:    config.listenAddr,
		Handler: snooper,
	}
	snooper.client = &http.Client{}
	return snooper, nil
}

func (s *Snooper) Snoop(ctx context.Context) error {
	go func() {
		if err := s.server.ListenAndServe(); err != nil {
			fmt.Printf("error: %s\n", err.Error())
		}
	}()
	for {
		<-ctx.Done()
		return s.server.Shutdown(ctx)
	}
}
