package registry

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

// TagData represents only necessary fields from maniest
type TagData struct {
	Name      string
	Version   string
	CreatedAt time.Time
}

// V1Compatibility represents a field from Manifest struct
type V1Compatibility struct {
	Architecture  string    `json:"architecture"`
	Created       time.Time `json:"created"`
	DockerVersion string    `json:"docker_version"`
	ID            string    `json:"id"`
	Os            string    `json:"os"`
}

// ImageData represents image object
type ImageData struct {
	Name    string
	Created time.Time
	Tag     string
}

// AllImages is used to get all the images
type AllImages struct {
	Images []ImageData
}

type RegistryClientOptions struct {
	Hostname string
}

type registryClient struct {
	hostname   string
	httpClient *http.Client
}

func New(options RegistryClientOptions) *registryClient {

	basicAuthTransport := &BasicTransport{
		Transport: http.DefaultTransport,
		URL:       options.Hostname,
		Username:  "",
		Password:  "",
	}
	errorTransport := &ErrorTransport{
		Transport: basicAuthTransport,
	}

	return &registryClient{
		hostname: options.Hostname,
		httpClient: &http.Client{
			Transport: errorTransport,
		},
	}
}

func (rc *registryClient) requestAndGetBody(query string) ([]byte, error) {

	req, err := http.NewRequest("GET", query, nil)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth("", "")

	resp, err := rc.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	bodyText, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return bodyText, nil
}

type BasicTransport struct {
	Transport http.RoundTripper
	URL       string
	Username  string
	Password  string
}

func (t *BasicTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if strings.HasPrefix(req.URL.String(), t.URL) {
		if t.Username != "" || t.Password != "" {
			req.SetBasicAuth(t.Username, t.Password)
		}
	}
	resp, err := t.Transport.RoundTrip(req)
	return resp, err
}

type HTTPStatusError struct {
	Response *http.Response
	// Copied from `Response.Body` to avoid problems with unclosed bodies later.
	// Nobody calls `err.Response.Body.Close()`, ever.
	Body []byte
}

func (err *HTTPStatusError) Error() string {
	return fmt.Sprintf("http: non-successful response (status=%v body=%q)", err.Response.StatusCode, err.Body)
}

var _ error = &HTTPStatusError{}

type ErrorTransport struct {
	Transport http.RoundTripper
}

func (t *ErrorTransport) RoundTrip(request *http.Request) (*http.Response, error) {
	resp, err := t.Transport.RoundTrip(request)
	if err != nil {
		return resp, err
	}

	if resp.StatusCode >= 400 {
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("http: failed to read response body (status=%v, err=%q)", resp.StatusCode, err)
		}

		return nil, &HTTPStatusError{
			Response: resp,
			Body:     body,
		}
	}

	return resp, err
}
