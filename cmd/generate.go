/*
Copyright © 2021 HugoByte <hello@hugobyte.com>

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

	"github.com/hugobyte/keygen/keystore/near"
	"github.com/spf13/cobra"
)

var password string

var out string

// generateCmd represents the generate command
var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generates KeyStore form private key ",
	Long:  "This Commmand Generates Keystore from the newly generated Publickey and Private Key pair",

	Run: func(cmd *cobra.Command, args []string) {

		err := near.GenerateNewKeystore(out, password)
		if err != nil {
			fmt.Println(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(generateCmd)

	generateCmd.Flags().StringVarP(&password, "pass", "p", "", "Password to Create KeyStore")
	generateCmd.Flags().StringVarP(&out, "out", "o", "keystore.json", "OutPut file path")
	generateCmd.MarkFlagRequired("pass")

}
