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
	{"swimlaneWorkstates", "s", "display the swimlane workstates of given squad"},
}

// Cmd represents the squad command
var Cmd = &cobra.Command{
	Use:           "squad",
	Short:         "display information for a given squad",
	SilenceErrors: true,
	Long: `display information for a given squad
	`,
	RunE: squadMain,
}

func squadMain(cmd *cobra.Command, args []string) error {

	vault.ReadVault()

	for _, f := range flags {
		if cmd.Flag(f.Name).Changed {
			var gQL graphQL
			switch f.Name {
			case "swimlaneWorkstates":
				if err := gQL.displayswimlaneWorkstates(args[0]); err != nil {
					return err
				}

			}

		}
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

	SwimlaneWorkstatesTotal := len(gQL.Data.Squad.SwimlaneWorkstates)
	var BacklogWorkstatesTotal int
	var ActiveWorkstatesTotal int
	var WaitWorkstatesTotal int

	for _, sl := range gQL.Data.Squad.SwimlaneWorkstates {
		fmt.Println(sl.Name)
		printTree("Backlog WorkStates", sl.BacklogWorkstates)
		printTree("Active WorkStates", sl.ActiveWorkstates)
		printTree("Wait WorkStates", sl.WaitWorkstates)
		fmt.Println()
		BacklogWorkstatesTotal += len(sl.BacklogWorkstates)
		ActiveWorkstatesTotal += len(sl.ActiveWorkstates)
		WaitWorkstatesTotal += len(sl.WaitWorkstates)

	}

	fmt.Printf("%v SwimlaneWorkstates, %v BacklogWorkstates, %v ActiveWorkstates,  %v WaitWorkstates\n",
		SwimlaneWorkstatesTotal, BacklogWorkstatesTotal, ActiveWorkstatesTotal, WaitWorkstatesTotal)
	return nil
}

func printTree(t string, s []string) {

	if len(s) == 0 {
		return
	}

	fmt.Printf("├── %s:\n", t)
	for i := 0; i < len(s); i++ {
		fmt.Println("│   └──", s[i])

	}

}

func init() {

	for _, f := range flags {
		Cmd.Flags().BoolP(f.Name, f.Short, false, f.Description)
	}
	Cmd.Flags().String("cardEmail", "", "grep cards with for given email (works with card only)")
	Cmd.Flags().String("cardSL", "", "grep cards with for given SwinLane (works with card only)")
	Cmd.Flags().String("cardPL", "", "grep cards with for given Primary Label (works with card only)")

	Cmd.PersistentFlags().String("id", "", "squad id")
	Cmd.AddCommand(usersCmd)
	Cmd.AddCommand(DescCmd)
	Cmd.AddCommand(cardsCmd)

	cardsCmd.Flags().String("filterEmail", "", "filter for cards for the email")
	cardsCmd.Flags().String("filterSLane", "", "filter for cards for the SwinLane")
	cardsCmd.Flags().String("filterPLabel", "", "filter for cards for the Primary Label")

	now := time.Now()
	lastMonth := now.AddDate(0, -1, 0)
	cardsCmd.Flags().String("updatedSince", lastMonth.Format(time.RFC3339), "filter for cards for the Primary Label")

	if err := Cmd.MarkPersistentFlagRequired("id"); err != nil {
		return
	}
	vault.ReadVault()
}
