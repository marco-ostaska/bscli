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

package vault

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

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

func addCommandUpdateCmd() error {
	Cmd.AddCommand(updateCmd)
	updateCmd.Flags().StringP("key", "k", "", "API key value")

	return updateCmd.MarkFlagRequired("key")

}

func updateVault(cmd *cobra.Command, args []string) error {
	err := Credential.ReadFile(Dir, File)
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

	fmt.Println("Updating credentials", Credential.URL)

	if err = Credential.SetInfo(APIKey, keyValue, Credential.URL, Dir, File); err != nil {
		return err
	}
	fmt.Println("Vault configured ✔")
	return nil

}
