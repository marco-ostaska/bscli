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

	"github.com/spf13/cobra"
)

// graphQL most primitive data for squad resturns
type graphQL struct {
	Data data `json:"data"`
}

// squad's user information
type users struct {
	Fullname string `json:"fullname"`
	Email    string `json:"email"`
}

// squad's Assignees information
type assignees struct {
	Email    string `json:"email"`
	Fullname string `json:"fullname"`
}

// squad's cards information
type cards struct {
	Identifier     string      `json:"identifier"`
	Title          string      `json:"title"`
	Description    string      `json:"description"`
	PrimaryLabels  []string    `json:"primaryLabels"`
	SecondaryLabel string      `json:"secondaryLabel"`
	DueAt          string      `json:"dueAt"`
	Swimlane       string      `json:"swimlane"`
	WorkstateType  string      `json:"workstateType"`
	Assignees      []assignees `json:"assignees"`
}

// squad's SwimlaneWorkstates information
type swimlaneWorkstates struct {
	Name string `json:"name"`
}

// Squad is an abstraction to squad
type squad struct {
	Name               string               `json:"name"`
	Users              []users              `json:"users"`
	Description        string               `json:"description"`
	Geography          string               `json:"geography"`
	SquadUsersCount    int                  `json:"squadUsersCount"`
	Cards              []cards              `json:"cards"`
	SwimlaneWorkstates []swimlaneWorkstates `json:"swimlaneWorkstates"`
}

// Data is the squad data
type data struct {
	Squad squad `json:"squad"`
}

// squadCmd represents the squad command
var squadCmd = &cobra.Command{
	Use:           "squad [id]...",
	Args:          cobra.ExactArgs(1),
	Short:         "display information for a given squad",
	SilenceErrors: true,
	Long: `display information for a given squad
	`,

	RunE: initSquad,
}

func initSquad(cmd *cobra.Command, args []string) error {
	if cmd.Flags().NFlag() < 1 {
		return fmt.Errorf("no flag(s) provided")
	}
	readVault()

	if err := displayName(cmd, args[0]); err != nil {
		log.Fatalln(err)
	}

	if err := displayDescription(cmd, args[0]); err != nil {
		log.Fatalln(err)
	}

	if err := displayUsers(cmd, args[0]); err != nil {
		return err
	}

	return nil

}

func displayUsers(cmd *cobra.Command, id string) error {
	u, err := cmd.Flags().GetBool("user")
	if err != nil {
		return err
	}

	if !u {
		return nil
	}

	var gQL graphQL
	query := fmt.Sprintf(`{
		squad(id: %s) {
		  users{
			fullname
			email
		  }
		  squadUsersCount
		}
	  }`, id)

	if err := httpc.QueryGraphQL(query, &gQL); err != nil {
		return err
	}

	fmt.Println("Users:")
	for _, u := range gQL.Data.Squad.Users {
		fmt.Printf("- %s (%s)\n", u.Fullname, u.Email)
	}
	fmt.Printf("\ntotal: %v\n", gQL.Data.Squad.SquadUsersCount)

	return nil
}

// TODO: need a better way of doing this, too much imilar code
func displayDescription(cmd *cobra.Command, id string) error {
	d, err := cmd.Flags().GetBool("description")
	if err != nil {
		return err
	}

	if !d {
		return nil
	}
	var gQL graphQL

	query := fmt.Sprintf(`{
		squad(id: %s) {
		  description
		}
	  }
	  `, id)

	if err := httpc.QueryGraphQL(query, &gQL); err != nil {
		return err
	}

	fmt.Println("Description:", gQL.Data.Squad.Description)

	return nil
}

func displayName(cmd *cobra.Command, id string) error {
	n, err := cmd.Flags().GetBool("name")
	if err != nil {
		return err
	}

	if !n {
		return nil
	}
	var gQL graphQL

	query := fmt.Sprintf(`{
		squad(id: %s) {
		  name
		}
	  }
	  `, id)

	if err := httpc.QueryGraphQL(query, &gQL); err != nil {
		return err
	}

	fmt.Println("Name:", gQL.Data.Squad.Name)

	return nil
}

func init() {
	rootCmd.AddCommand(squadCmd)
	squadCmd.Flags().BoolP("name", "n", false, "display name for given squad")
	squadCmd.Flags().BoolP("user", "u", false, "display users for given squad")
	squadCmd.Flags().BoolP("description", "d", false, "display description for given squad")

}
