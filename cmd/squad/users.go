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

// Package squad references bs1cli squad command used by cobra.
package squad

import (
	"fmt"

	"github.com/marco-ostaska/bscli/cmd/vault"
	"github.com/spf13/cobra"
)

// usersCmd represents the squad command
var usersCmd = &cobra.Command{
	Use:           "users",
	Short:         "display the users for the squad",
	SilenceErrors: true,
	Example:       `bscli squad --id <squad id> users`,
	Long: `display the users for the squad
	`,
	RunE: displayUsers,
}

func displayUsers(cmd *cobra.Command, args []string) error {
	vault.ReadVault()

	id, err := cmd.Flags().GetString("id")
	if err != nil {
		return err
	}

	var gQL graphQL
	query := fmt.Sprintf(`{
		squad(id: %s) {
		  users{
			email
			fullname
		  }
		  squadUsersCount
		}
	  }`, id)

	if err := vault.HTTP.QueryGraphQL(query, &gQL); err != nil {
		return err
	}

	for _, u := range gQL.Data.Squad.Users {
		fmt.Printf("- %s (%s)\n", u.Fullname, u.Email)
	}

	return nil
}

func init() {
	Cmd.AddCommand(usersCmd)
}
