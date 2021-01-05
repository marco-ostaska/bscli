/*
Copyright © 2021 Marco Ostaska

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.
*/

package cmd

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/marco-ostaska/boringstuff"
	"github.com/marco-ostaska/uvault"
	"github.com/spf13/cobra"
)

const (
	vaultDir  = "bscli"
	vaultFile = "bscli.vlt"
	apiKey    = "Bluesight-API-Token"
)

var vCredential uvault.Credential

// vaultCmd represents the vault command
var vaultCmd = &cobra.Command{
	Use:   "vault",
	Short: "create or update vault credentials",
	Long:  `create or update vault credentials.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("%v\n", cmd.Long)
		if err := cmd.Usage(); err != nil {
			log.Fatalln(err)
		}
	},
}

var newCmd = &cobra.Command{
	Use:           "new",
	Short:         "create new vault.",
	Long:          `create new vault.`,
	SilenceErrors: true,
	Example: `
  Unix Based OS: (use single quotes)
      sl1cmd vault new -k 'pass1234' --url 'https://bluesight.com'
  Windows: (use double quotes)
      sl1cmd vault new -k "pass1234" --url "https://bluesight.com"
`,
	RunE: newVault,
}

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "update an existing vault.",
	Long:  `update an existing vault.`,
	Example: `  
  Unix based OS:  (use single quotes)
      sl1cmd update -k 'pass1234'
  Windows: (use double quotes)
      sl1cmd update -k "pass1234"`,
	RunE: updateVault,
}

var deleteCmd = &cobra.Command{
	Use:           "delete",
	Short:         "delete an existing vault.",
	Long:          `delete an existing vault.`,
	SilenceUsage:  true,
	SilenceErrors: true,
	RunE:          deleteVault,
}

func addCommandUpdateCmd() error {
	vaultCmd.AddCommand(updateCmd)
	updateCmd.Flags().StringP("key", "k", "", "API key value")

	return updateCmd.MarkFlagRequired("key")

}

func addCommandNewCmd() error {
	vaultCmd.AddCommand(newCmd)
	newCmd.Flags().StringP("key", "k", "", "API key value")
	newCmd.Flags().String("url", "", "API URI")

	err := newCmd.MarkFlagRequired("key")
	err1 := newCmd.MarkFlagRequired("url")

	if re := boringstuff.ReturnError(err, err1); re != nil {
		return re
	}

	return nil
}

func newVault(cmd *cobra.Command, args []string) error {
	keyValue, err := cmd.Flags().GetString("key")
	uri, err1 := cmd.Flags().GetString("url")

	if re := boringstuff.ReturnError(err, err1); re != nil {
		return re
	}

	if err = vCredential.SetInfo(apiKey, keyValue, uri, vaultDir, vaultFile); err != nil {
		return err
	}

	fmt.Println("Vault configured ✔")
	return nil
}

func updateVault(cmd *cobra.Command, args []string) error {
	err := vCredential.ReadFile(vaultDir, vaultFile)
	if err != nil {
		if strings.Contains(err.Error(), "no such file or directory") {
			return fmt.Errorf("No credentials found, please try create a new credential vault first ❌")
		}
		return err
	}
	keyValue, err := cmd.Flags().GetString("key")
	if err != nil {
		return err
	}

	fmt.Println("Updating credentials", vCredential.URL)

	if err = vCredential.SetInfo(apiKey, keyValue, vCredential.URL, vaultDir, vaultFile); err != nil {
		return err
	}
	fmt.Println("Vault configured ✔")
	return nil

}

func deleteVault(cmd *cobra.Command, args []string) error {
	if err := vCredential.UserInfo(vaultDir, vaultFile); err != nil {
		return fmt.Errorf("%s ❌", err)
	}

	if err := os.Remove(vCredential.File); err != nil {
		return fmt.Errorf("%s ❌", err.Error())
	}

	fmt.Println("Vault deleted ✔")
	return nil

}

func init() {
	rootCmd.AddCommand(vaultCmd)
	vaultCmd.AddCommand(deleteCmd)

	err := addCommandNewCmd()
	err1 := addCommandUpdateCmd()

	if re := boringstuff.ReturnError(err, err1); re != nil {
		log.Fatalln(err)
	}

}
