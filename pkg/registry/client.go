package registry

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/docker/distribution/manifest/schema1"

	//"strings"
	"time"
)

type repositoriesResponse struct {
	Repositories []string `json:"repositories"`
}

type tagsResponse struct {
	Name string   `json:"name"`
	Tags []string `json:"tags"`
}

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

func (rc *registryClient) GetImageData(repo string) ([]ImageData, error) {

	bodyText, err := rc.requestAndGetBody(fmt.Sprintf("https://%s/v2/%s/tags/list", rc.hostname, repo))
	if err != nil {
		return nil, err
	}

	tagsResp := tagsResponse{}

	err = json.Unmarshal(bodyText, &tagsResp)
	if err != nil {
		return nil, err
	}

	returnImageData := []ImageData{}
	for _, tag := range tagsResp.Tags {
		bodyText1, err := rc.requestAndGetBody(fmt.Sprintf("https://%s/v2/%s/manifests/%s", rc.hostname, tagsResp.Name, tag))
		if err != nil {
			return nil, err
		}

		mani := schema1.Manifest{}

		err = json.Unmarshal(bodyText1, &mani)
		if err != nil {
			return nil, err
		}

		v1Compatibility := V1Compatibility{}

		err = json.Unmarshal([]byte(mani.History[0].V1Compatibility), &v1Compatibility)
		if err != nil {
			return nil, err
		}

		returnImageData = append(returnImageData, ImageData{
			Name:    tagsResp.Name,
			Tag:     tag,
			Created: v1Compatibility.Created,
		})

	}

	return returnImageData, nil

}

func (rc *registryClient) GetRepos() ([]ImageData, error) {
	bodyText, err := rc.requestAndGetBody(fmt.Sprintf("https://%s/v2/_catalog", rc.hostname))
	if err != nil {
		return nil, err
	}

	repoResp := repositoriesResponse{}
	err = json.Unmarshal(bodyText, &repoResp)
	if err != nil {
		return nil, err
	}

	returnImageData := []ImageData{}
	for _, repo := range repoResp.Repositories {

		images, err := rc.GetImageData(repo)
		if err != nil {
			continue //return nil, err
		}

		returnImageData = append(returnImageData, images...)
	}

	return returnImageData, nil
}

func (rc *registryClient) requestAndGetBody(query string) ([]byte, error) {

	resp, err := rc.httpClient.Get(query)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

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
	if strings.Contains(req.URL.String(), t.URL) {
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
