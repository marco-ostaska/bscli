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
	"os"

	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:           "delete",
	Short:         "delete an existing vault.",
	Long:          `delete an existing vault.`,
	SilenceUsage:  true,
	SilenceErrors: true,
	RunE:          deleteVault,
}

func deleteVault(cmd *cobra.Command, args []string) error {
	if err := Credential.UserInfo(Dir, File); err != nil {
		return fmt.Errorf("%s ❌", err)
	}

	if err := os.Remove(Credential.File); err != nil {
		return fmt.Errorf("%s ❌", err.Error())
	}

	fmt.Println("Vault deleted ✔")
	return nil

}
