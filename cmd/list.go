/*

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/docker/distribution/manifest/schema1"
	"github.com/spf13/cobra"
)

type listOptions struct {
	hostname string
}

var listOpts = &listOptions{}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "list a container registry",
	Long:  "list a container registry",
	Example: "	bupkis list",
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		listOpts.hostname = args[0]
		return runList()
	},
}

func init() {
	listCmd.Flags().StringVarP(&listOpts.hostname, "hostname", "n", "", "registry hostname")
	RootCmd.AddCommand(listCmd)
}

type teer struct {
	Repositories []string `json:"repositories"`
}

type gett struct {
	Name string   `json:"name"`
	Tags []string `json:"tags"`
}

type hmm struct {
	Created string `json:"created"`
}

func runList() error {

	bodyText, err := requestAndGetBody(fmt.Sprintf("https://%s/v2/_catalog", listOpts.hostname))
	if err != nil {
		return err
	}

	test := teer{}

	err = json.Unmarshal(bodyText, &test)
	if err != nil {
		return err
	}

	for _, asf := range test.Repositories {

		bodyText, err := requestAndGetBody(fmt.Sprintf("https://%s/v2/%s/tags/list", listOpts.hostname, asf))
		if err != nil {
			return err
		}

		test := gett{}

		err = json.Unmarshal(bodyText, &test)
		if err != nil {
			return err
		}

		fmt.Printf("%s %s\n", test.Name, test.Tags)

		bodyText1, err := requestAndGetBody(fmt.Sprintf("https://%s/v2/%s/manifests/%s", listOpts.hostname, test.Name, test.Tags[0]))
		if err != nil {
			return err
		}

		test1 := schema1.Manifest{}

		err = json.Unmarshal(bodyText1, &test1)
		if err != nil {
			return err
		}

		fmt.Println() //(\n", test1.History[0])

		test2 := hmm{}

		err = json.Unmarshal([]byte(test1.History[0].V1Compatibility), &test2)
		if err != nil {
			return err
		}
		fmt.Printf("%s\n", test2.Created)

	}

	return nil
}

func requestAndGetBody(query string) ([]byte, error) {

	client := &http.Client{}

	req, err := http.NewRequest("GET", query, nil)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth("", "")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	bodyText, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return bodyText, nil
}
