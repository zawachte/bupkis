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
	"fmt"
	"os"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"github.com/zawachte-msft/bupkis/pkg/registry"
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

func runList() error {

	client := registry.New(registry.RegistryClientOptions{Hostname: listOpts.hostname})

	images, err := client.GetRepos()
	if err != nil {
		return err
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Name", "Tag", "CreatedAt"})
	table.SetAutoWrapText(false)
	table.SetAutoFormatHeaders(true)
	table.SetHeaderAlignment(ALIGN_LEFT)
	table.SetAlignment(ALIGN_LEFT)
	table.SetCenterSeparator("")
	table.SetColumnSeparator("")
	table.SetRowSeparator("")
	table.SetHeaderLine(false)
	table.SetBorder(false)
	table.SetTablePadding("\t") // pad with tabs
	table.SetNoWhiteSpace(true)

	imageData := [][]string{}
	for _, image := range images {

		imageName := fmt.Sprintf("%s/%s", listOpts.hostname, image.Name)
		//fmt.Printf("%s\t%s\n", imageName, image.Created)
		table.Append([]string{imageName, image.Tag, image.Created.Format("2 Jan 2006 15:04:05")})
	}
	table.Render()

	return nil
}
