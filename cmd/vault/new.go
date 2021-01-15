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

	"github.com/marco-ostaska/boringstuff"
	"github.com/spf13/cobra"
)

var newCmd = &cobra.Command{
	Use:   "new",
	Short: "create new vault.",
	Long: `create new vault.
For more information how to create a token please 
check the API section https://portal.bluesight.io/tutorial.html 
`,
	SilenceErrors: true,
	Example: `
  Unix Based OS: (use single quotes)
      bscli vault new -k '<token>' --url 'https://www.bluesight.io/graphql'
  Windows: (use double quotes)
      bscli vault new -k "<token>" --url "https://www.bluesight.io/graphql"
`,
	RunE: newVault,
}

func addCommandNewCmd() error {
	Cmd.AddCommand(newCmd)
	newCmd.Flags().StringP("key", "k", "", "API key value")
	newCmd.Flags().String("url", "", "API URI")

	errs := []error{
		newCmd.MarkFlagRequired("key"),
		newCmd.MarkFlagRequired("url"),
	}

	for _, err := range errs {
		if err != nil {
			return err
		}
	}

	return nil
}

func newVault(cmd *cobra.Command, args []string) error {

	keyValue, err := cmd.Flags().GetString("key")
	uri, err1 := cmd.Flags().GetString("url")

	if re := boringstuff.ReturnError(err, err1); re != nil {
		return re
	}

	if err = Credential.SetInfo(APIKey, keyValue, uri, Dir, File); err != nil {
		return err
	}

	fmt.Println("Vault configured ✔")
	return nil
}
