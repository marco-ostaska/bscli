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
	"time"

	"github.com/marco-ostaska/bscli/cmd/vault"
	"github.com/spf13/cobra"
)

// graphQL most primitive data for squad resturns
type graphQL struct {
	Data struct {
		Squad struct {
			ID          string `json:"id"`
			Description string `json:"description"`
			Name        string `json:"name"`
			Users       []struct {
				Email    string `json:"email"`
				Fullname string `json:"fullname"`
			} `json:"users"`
			SwimlaneWorkstates []struct {
				Name              string   `json:"name"`
				BacklogWorkstates []string `json:"backlogWorkstates"`
				ActiveWorkstates  []string `json:"activeWorkstates"`
				WaitWorkstates    []string `json:"waitWorkstates"`
			} `json:"swimlaneWorkstates"`
			Cards []struct {
				Identifier     string      `json:"identifier"`
				Title          string      `json:"title"`
				PrimaryLabels  []string    `json:"primaryLabels"`
				SecondaryLabel interface{} `json:"secondaryLabel"`
				Swimlane       string      `json:"swimlane"`
				WorkstateType  string      `json:"workstateType"`
				DueAt          string      `json:"dueAt"`
				Assignees      []struct {
					Fullname string `json:"fullname"`
					Email    string `json:"email"`
				} `json:"assignees"`
			} `json:"cards"`
		} `json:"squad"`
	} `json:"data"`
}

var flags = []struct {
	Name        string
	Short       string
	Description string
}{
	{"name", "n", "display the name of given squad"},
	{"users", "u", "display the users of given squad"},
	{"swimlaneWorkstates", "s", "display the swimlane workstates of given squad"},
	{"description", "d", "display the description of given squad"},
	{"cards", "c", "diplay active cards of the given squad (Only cards updated in the last month)"},
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
			var gQL graphQL
			switch f.Name {
			case "users":
				if err := gQL.displayUsers(args[0]); err != nil {
					return err
				}
			case "name":
				if err := gQL.displayName(args[0]); err != nil {
					return err
				}
			case "description":
				if err := gQL.displayDescription(args[0]); err != nil {
					return err
				}
			case "swimlaneWorkstates":
				if err := gQL.displayswimlaneWorkstates(args[0]); err != nil {
					return err
				}

			case "cards":
				email, err := cmd.Flags().GetString("cardEmail")
				if err != nil {
					return err
				}

				sl, err := cmd.Flags().GetString("cardSL")
				if err != nil {
					return err
				}
				pl, err := cmd.Flags().GetString("cardPL")
				if err != nil {
					return err
				}

				if err := gQL.displayCards(args[0], email, sl, pl); err != nil {
					return err
				}
			}

		}
	}
	return nil

}

func (gQL graphQL) displayUsers(id string) error {
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

func (gQL graphQL) displayswimlaneWorkstates(id string) error {
	query := fmt.Sprintf(`{
		squad(id: %s) {
		  swimlaneWorkstates{
			name
			backlogWorkstates
			activeWorkstates
			waitWorkstates
		  }
		}
	  }`, id)

	if err := vault.HTTP.QueryGraphQL(query, &gQL); err != nil {
		return err
	}

	for _, sl := range gQL.Data.Squad.SwimlaneWorkstates {
		fmt.Println(sl.Name)
		printTree("Backlog WorkStates", sl.BacklogWorkstates)
		printTree("Active WorkStates", sl.ActiveWorkstates)
		printTree("Wait WorkStates", sl.WaitWorkstates)
		fmt.Println()
	}
	return nil
}

func printTree(t string, s []string) {

	if len(s) == 0 {
		return
	}

	fmt.Printf("├── %s:\n", t)
	for i := 0; i < len(s); i++ {
		fmt.Println("|   └──", s[i])
	}

}

func (gQL graphQL) displayDescription(id string) error {
	query := fmt.Sprintf(`{
		squad(id: %s) {
		  description
		}
	  }
	  `, id)

	if err := vault.HTTP.QueryGraphQL(query, &gQL); err != nil {
		return err
	}
	fmt.Println(gQL.Data.Squad.Description)

	return nil
}

func (gQL graphQL) displayCards(id, email, sl, pl string) error {
	now := time.Now()
	lastMonth := now.AddDate(0, -1, 0)

	query := fmt.Sprintf(`{
	squad(id: %s){
		cards(updatedSince: "%v", closed: false){
		identifier
		title
		primaryLabels
		secondaryLabel
		swimlane
		workstateType
		dueAt
		assignees{
			fullname
			email
		}   
		}
	}
	}`, id, lastMonth.Format(time.RFC3339))
	if err := vault.HTTP.QueryGraphQL(query, &gQL); err != nil {
		return err
	}

	for _, c := range gQL.Data.Squad.Cards {

		var emailOK bool
		var slOK bool
		var plOK bool

		for _, a := range c.Assignees {
			if a.Email == email {
				emailOK = true
			}

		}

		for _, p := range c.PrimaryLabels {
			if p == pl {
				plOK = true
			}

		}

		if c.Swimlane == sl {
			slOK = true
		}

		switch {
		case len(email) == 0:
			emailOK = true
		case len(sl) == 0:
			slOK = true
		case len(pl) == 0:
			plOK = true
		}

		if !emailOK || !slOK || !plOK {
			continue
		}

		fmt.Println()
		fmt.Println("Identifier      :", c.Identifier)
		fmt.Println("Title           :", c.Title)
		fmt.Println("Work State      :", c.WorkstateType)
		fmt.Println("SwinLane        :", c.Swimlane)
		fmt.Println("Due Date        :", c.DueAt)
		fmt.Print("Assinee(s)      : [")
		for _, a := range c.Assignees {
			fmt.Printf(" %v ", a.Email)
		}
		fmt.Println("]")
		if len(c.PrimaryLabels) > 0 {
			fmt.Println("Primary Label(s):", c.PrimaryLabels)
		}

		if c.SecondaryLabel != nil {
			fmt.Println("Secondary Label :", c.SecondaryLabel)
		}

	}
	return nil
}

func (gQL graphQL) displayName(id string) error {
	query := fmt.Sprintf(`{
		squad(id: %s) {
		  name
		}
	  }
	  `, id)
	if err := vault.HTTP.QueryGraphQL(query, &gQL); err != nil {
		return err
	}
	fmt.Println(gQL.Data.Squad.Name)
	return nil
}

func init() {

	for _, f := range flags {
		Cmd.Flags().BoolP(f.Name, f.Short, false, f.Description)
	}

	Cmd.Flags().String("cardEmail", "", "grep cards with for given email (works with card only)")
	Cmd.Flags().String("cardSL", "", "grep cards with for given SwinLane (works with card only)")
	Cmd.Flags().String("cardPL", "", "grep cards with for given Primary Label (works with card only)")
}
