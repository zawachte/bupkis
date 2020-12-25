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
	"github.com/spf13/cobra"
	"github.com/zawachte-msft/bupkis/pkg/formatter"
	"github.com/zawachte-msft/bupkis/pkg/registry"
	"github.com/zawachte-msft/bupkis/pkg/util"
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
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 0 {
			listOpts.hostname = args[0]
		}

		return runList()
	},
}

func init() {
	listCmd.Flags().StringVarP(&listOpts.hostname, "hostname", "n", "", "registry hostname")
	RootCmd.AddCommand(listCmd)
}

func runList() error {

	client, err := registry.New(registry.RegistryClientOptions{Hostname: listOpts.hostname})
	if err != nil {
		return err
	}

	images, err := client.GetRepos()
	if err != nil {
		return err
	}

	formatter.PrintOutput(util.ImagesToNestedArray(images))

	return nil
}
