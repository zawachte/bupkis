package registry

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/docker/distribution/manifest/schema1"
	auth "github.com/zawachte-msft/bupkis/pkg/auth/docker"

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
	Name     string
	Created  time.Time
	Tag      string
	Hostname string
}

// AllImages is used to get all the images
type AllImages struct {
	Images []ImageData
}

type RegistryClientOptions struct {
	Hostname string
}

type registryClient struct {
	hostname      string
	httpClientMap map[string]*http.Client
}

func New(options RegistryClientOptions) (*registryClient, error) {

	httpClientMap := make(map[string]*http.Client)

	// Prepare auth client
	cli, err := auth.NewClient()
	if err != nil {
		return nil, err
	}

	if options.Hostname != "" {
		username, password, err := cli.Credential(options.Hostname)
		if err != nil {
			return nil, err
		}

		basicAuthTransport := &BasicTransport{
			Transport: http.DefaultTransport,
			URL:       options.Hostname,
			Username:  username,
			Password:  password,
		}
		errorTransport := &ErrorTransport{
			Transport: basicAuthTransport,
		}

		httpClientMap[options.Hostname] = &http.Client{
			Transport: errorTransport,
		}
	} else {
		authConfigMap, err := cli.GetAllCredentials()
		if err != nil {
			return nil, err
		}

		for hostname, authConfig := range authConfigMap {
			basicAuthTransport := &BasicTransport{
				Transport: http.DefaultTransport,
				URL:       hostname,
				Username:  authConfig.Username,
				Password:  authConfig.Password,
			}
			errorTransport := &ErrorTransport{
				Transport: basicAuthTransport,
			}

			httpClientMap[hostname] = &http.Client{
				Transport: errorTransport,
			}
		}
	}

	return &registryClient{
		hostname:      options.Hostname,
		httpClientMap: httpClientMap,
	}, nil
}

func (rc *registryClient) GetImageDataList(hostname string, repo string) ([]ImageData, error) {

	bodyText, err := rc.requestAndGetBody(hostname, fmt.Sprintf("https://%s/v2/%s/tags/list", hostname, repo))
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

		imageData, err := rc.GetImageData(hostname, tagsResp.Name, tag)
		if err != nil {
			return nil, err
		}

		returnImageData = append(returnImageData, imageData)
	}

	return returnImageData, nil

}
func (rc *registryClient) GetImageData(hostname string, repo string, tag string) (ImageData, error) {
	bodyText1, err := rc.requestAndGetBody(hostname, fmt.Sprintf("https://%s/v2/%s/manifests/%s", hostname, repo, tag))
	if err != nil {
		return ImageData{}, err
	}

	mani := schema1.Manifest{}

	err = json.Unmarshal(bodyText1, &mani)
	if err != nil {
		return ImageData{}, err
	}

	v1Compatibility := V1Compatibility{}

	err = json.Unmarshal([]byte(mani.History[0].V1Compatibility), &v1Compatibility)
	if err != nil {
		return ImageData{}, err
	}

	return ImageData{
		Name:     repo,
		Tag:      tag,
		Created:  v1Compatibility.Created,
		Hostname: hostname,
	}, nil

}

func (rc *registryClient) GetRepos() ([]ImageData, error) {

	returnImageData := []ImageData{}

	for hostname := range rc.httpClientMap {
		imageData, err := rc.GetReposByHostName(hostname)
		if err != nil {
			return nil, err
		}

		returnImageData = append(returnImageData, imageData...)
	}

	return returnImageData, nil
}

func (rc *registryClient) GetReposByHostName(hostname string) ([]ImageData, error) {

	returnImageData := []ImageData{}
	bodyText, err := rc.requestAndGetBody(hostname, fmt.Sprintf("https://%s/v2/_catalog", hostname))
	if err != nil {
		return nil, err
	}

	repoResp := repositoriesResponse{}
	err = json.Unmarshal(bodyText, &repoResp)
	if err != nil {
		return nil, err
	}

	for _, repo := range repoResp.Repositories {

		images, err := rc.GetImageDataList(hostname, repo)
		if err != nil {
			continue
		}

		returnImageData = append(returnImageData, images...)
	}

	return returnImageData, nil
}

func (rc *registryClient) requestAndGetBody(hostname string, query string) ([]byte, error) {

	resp, err := rc.httpClientMap[hostname].Get(query)
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
