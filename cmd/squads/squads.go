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

package squads

import (
	"fmt"

	"github.com/marco-ostaska/bscli/cmd/vault"
	"github.com/spf13/cobra"
)

// Cmd represents the squads command
var Cmd = &cobra.Command{
	Use:           "squads",
	Short:         "list the squads for the user",
	Long:          `list the squads for the user`,
	SilenceErrors: true,
	SilenceUsage:  true,
	RunE:          squadsDisplay,
}

type graphQL struct {
	Data struct {
		Squads []struct {
			ID   string `json:"id"`
			Name string `json:"name"`
		} `json:"squads"`
	} `json:"data"`
}

func squadsDisplay(cmd *cobra.Command, args []string) error {
	vault.ReadVault()

	var gQL graphQL
	query := `{
		squads {
		  id
		  name
		  squadUsersCount
		}
	  }
	  `

	if err := vault.HTTP.QueryGraphQL(query, &gQL); err != nil {
		if err.Error() == vault.ErrLoginFailure.Error() {
			return fmt.Errorf("Login Failure, please check your token and url and try again")
		}
		return err
	}

	for _, v := range gQL.Data.Squads {
		fmt.Printf("id=%s(%s)\n", v.ID, v.Name)
	}

	return nil
}
