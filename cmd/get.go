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

type getOptions struct {
	image string
}

var getOpts = &getOptions{}

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "get a container registry",
	Long:  "get a container registry",
	Example: "	bupkis get",
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		getOpts.image = args[0]
		return runGet()
	},
}

func init() {
	RootCmd.AddCommand(getCmd)
}

func runGet() error {

	imageData := util.ParseImageName(getOpts.image)

	client, err := registry.New(registry.RegistryClientOptions{Hostname: imageData.Hostname})
	if err != nil {
		return err
	}

	imagesDatas := []registry.ImageData{}

	if imageData.Tag == "" {
		images, err := client.GetImageDataList(imageData.Hostname, imageData.Name)
		if err != nil {
			return err
		}

		imagesDatas = append(imagesDatas, images...)
	} else {
		images, err := client.GetImageData(imageData.Hostname, imageData.Name, imageData.Tag)
		if err != nil {
			return err
		}

		imagesDatas = append(imagesDatas, images)
	}

	formatter.PrintOutput(util.ImagesToNestedArray(imagesDatas))

	return nil
}
