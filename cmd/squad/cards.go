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

package squad

import (
	"fmt"
	"time"

	"github.com/marco-ostaska/bscli/cmd/vault"
	"github.com/spf13/cobra"
)

// cardsCmd represents the squad command
var cardsCmd = &cobra.Command{
	Use:           "cards",
	Short:         "diplay active cards of the given squad (Only cards updated in the last month)",
	SilenceErrors: true,
	Example:       `bscli squad --id <squad id> cards --filterEmail my@email.com --updatedSince "2020-01-31T22:37:22-03:00" `,
	Long: `display the users for the squad
	`,
	RunE: displayCards,
}

type assignees []struct {
	Fullname string `json:"fullname"`
	Email    string `json:"email"`
}

func filter(cmd cobra.Command, asg assignees, pLabel []string, pSlane string) bool {

	filterEmail, err := cmd.Flags().GetString("filterEmail")
	if err != nil {
		return false
	}
	filterSlane, err := cmd.Flags().GetString("filterSLane")
	if err != nil {
		return false
	}
	filterPlabel, err := cmd.Flags().GetString("filterPLabel")
	if err != nil {
		return false
	}

	var emailOK bool
	var slOK bool
	var plOK bool

	for _, a := range asg {
		if a.Email == filterEmail {
			emailOK = true
		}
	}

	for _, p := range pLabel {
		if p == filterPlabel {
			plOK = true
		}
	}

	if pSlane == filterSlane {
		slOK = true
	}

	if len(filterEmail) == 0 {
		emailOK = true
	}

	if len(filterSlane) == 0 {
		slOK = true
	}

	if len(filterPlabel) == 0 {
		plOK = true
	}

	if !emailOK || !slOK || !plOK {
		return false
	}
	return true
}

func displayCards(cmd *cobra.Command, args []string) error {
	vault.ReadVault()

	updatedSince, err := cmd.Flags().GetString("updatedSince")
	if err != nil {
		return err
	}

	id, err := cmd.Flags().GetString("id")
	if err != nil {
		return err
	}

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
	}`, id, updatedSince)

	var gQL graphQL

	if err := vault.HTTP.QueryGraphQL(query, &gQL); err != nil {
		return err
	}

	for _, c := range gQL.Data.Squad.Cards {

		if !filter(*cmd, c.Assignees, c.PrimaryLabels, c.Swimlane) {
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

func init() {

	// cards commands and routines
	Cmd.AddCommand(cardsCmd)
	cardsCmd.Flags().String("filterEmail", "", "filter for cards for the email")
	cardsCmd.Flags().String("filterSLane", "", "filter for cards for the SwimLane")
	cardsCmd.Flags().String("filterPLabel", "", "filter for cards for the Primary Label")
	now := time.Now()
	lastMonth := now.AddDate(0, -1, 0)
	cardsCmd.Flags().String("updatedSince", lastMonth.Format(time.RFC3339), "filter for cards for the Primary Label")
}
