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
	"context"
	"fmt"

	"github.com/spf13/cobra"
	auth "github.com/zawachte-msft/bupkis/pkg/auth/docker"
)

type loginOptions struct {
	hostname  string
	manifest  string
	username  string
	password  string
	fromStdin bool
}

var loginOpts = &loginOptions{}

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "login into a container registry",
	Long:  "login into a container registry",
	Example: "	bupkis login cr.io --username bupkisuser --password-stdin",
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		loginOpts.hostname = args[0]
		return runLogin()
	},
}

func init() {
	loginCmd.Flags().StringVarP(&loginOpts.username, "username", "u", "", "registry username")
	loginCmd.Flags().StringVarP(&loginOpts.password, "password", "p", "", "registry password or identity token")
	loginCmd.Flags().BoolVarP(&loginOpts.fromStdin, "password-stdin", "", false, "read password or identity token from stdin")
	RootCmd.AddCommand(loginCmd)
}

func runLogin() error {

	// Prepare auth client
	cli, err := auth.NewClient()
	if err != nil {
		return err
	}

	// Login
	if err := cli.Login(context.Background(), loginOpts.hostname, loginOpts.username, loginOpts.password, false); err != nil {
		return err
	}
	fmt.Println("Login Succeeded")
	return nil
}
