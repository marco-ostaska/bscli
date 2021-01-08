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

var flags = []struct {
	Name  string
	Short string
	Call  func(string) error
}{
	{"name", "n", displayName},
	{"users", "u", displayUsers},
	{"swimlaneWorkstates", "s", displayswimlaneWorkstates},
	{"description", "d", displayDescription},
}

// Cmd represents the squad command
var Cmd = &cobra.Command{
	Use:           "squad [id]...",
	Args:          cobra.ExactArgs(1),
	Short:         "display information for a given squad",
	SilenceErrors: true,
	Long: `display information for a given squad
	`,

	RunE: squadMain,
}

func squadMain(cmd *cobra.Command, args []string) error {

	if cmd.Flags().NFlag() < 1 {
		return fmt.Errorf("No flag(s) received")
	}
	vault.ReadVault()

	for _, f := range flags {
		if cmd.Flag(f.Name).Changed {
			f.Call(args[0])
		}
	}
	return nil

}

func queryQL(query, flag string) error {
	var gQL graphQL

	if err := vault.HTTP.QueryGraphQL(query, &gQL); err != nil {
		return err
	}

	switch flag {
	case "users":
		for _, u := range gQL.Data.Squad.Users {
			fmt.Printf("- %s (%s)\n", u.Fullname, u.Email)
		}
		fmt.Printf("\ntotal: %v\n", gQL.Data.Squad.SquadUsersCount)
	case "name":
		fmt.Println(gQL.Data.Squad.Name)

	case "description":
		fmt.Println(gQL.Data.Squad.Description)
	case "swimlaneWorkstates":
		for _, u := range gQL.Data.Squad.SwimlaneWorkstates {
			fmt.Printf("- %s\n", u.Name)
		}
		fmt.Printf("\ntotal: %v\n", len(gQL.Data.Squad.SwimlaneWorkstates))

	default:
		return fmt.Errorf("unknown flag")

	}
	return nil
}

func displayUsers(id string) error {

	query := fmt.Sprintf(`{
		squad(id: %s) {
		  users{
			email
			fullname
		  }
		  squadUsersCount
		}
	  }`, id)

	return queryQL(query, "users")

}

func displayswimlaneWorkstates(id string) error {
	query := fmt.Sprintf(`{
		squad(id: %s) {
		  swimlaneWorkstates{
			activeWorkstates
			name
		  }
		}
	  }`, id)

	return queryQL(query, "swimlaneWorkstates")
}

// TODO: need a better way of doing this
func displayDescription(id string) error {
	query := fmt.Sprintf(`{
		squad(id: %s) {
		  description
		}
	  }
	  `, id)

	return queryQL(query, "description")
}

func displayName(id string) error {
	query := fmt.Sprintf(`{
		squad(id: %s) {
		  name
		}
	  }
	  `, id)
	return queryQL(query, "name")
}

func init() {

	for _, f := range flags {
		desc := fmt.Sprintf("display %s for given squad", f.Name)
		Cmd.Flags().BoolP(f.Name, f.Short, false, desc)
	}

}
