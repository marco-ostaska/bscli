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

type filters struct {
	filterEmail struct {
		value string
		ok    bool
	}
	filterSlane struct {
		value string
		ok    bool
	}
	filterPlabel struct {
		value string
		ok    bool
	}
	filterWorkState struct {
		value string
		ok    bool
	}
	assignees []struct {
		Fullname string `json:"fullname"`
		Email    string `json:"email"`
	}
	plabel    []string
	workState string
	pslane    string
}

func (f *filters) filterLen() {
	if len(f.filterEmail.value) == 0 {
		f.filterEmail.ok = true
	}

	if len(f.filterPlabel.value) == 0 {
		f.filterPlabel.ok = true
	}

	if len(f.filterSlane.value) == 0 {
		f.filterSlane.ok = true
	}

	if len(f.filterWorkState.value) == 0 {
		f.filterWorkState.ok = true
	}

}

func (f *filters) hasFilterOn(cmd cobra.Command) error {
	var errs [4]error

	f.filterEmail.value, errs[0] = cmd.Flags().GetString("filterEmail")
	f.filterSlane.value, errs[1] = cmd.Flags().GetString("filterSLane")
	f.filterPlabel.value, errs[2] = cmd.Flags().GetString("filterPLabel")
	f.filterWorkState.value, errs[3] = cmd.Flags().GetString("filterWorkState")

	for _, err := range errs {
		if err != nil {
			return err
		}
	}

	f.filterLen()

	return nil

}

//TODO break this function
func (f *filters) filter(cmd cobra.Command) bool {

	if err := f.hasFilterOn(cmd); err != nil {
		return false
	}

	for _, a := range f.assignees {
		if a.Email == f.filterEmail.value {
			f.filterEmail.ok = true
		}
	}

	for _, p := range f.plabel {
		if p == f.filterPlabel.value {
			f.filterPlabel.ok = true
		}
	}

	if f.workState == f.filterWorkState.value {
		f.filterWorkState.ok = true
	}

	if f.pslane == f.filterSlane.value {
		f.filterSlane.ok = true
	}

	if !f.filterSlane.ok || !f.filterPlabel.ok || !f.filterSlane.ok || !f.filterWorkState.ok {
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

		var f filters
		f.plabel = c.PrimaryLabels
		f.assignees = c.Assignees
		f.pslane = c.Swimlane
		f.workState = c.WorkstateType

		if !f.filter(*cmd) {
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
	cardsCmd.Flags().String("filterWorkState", "", "filter for cards for the WorkState Type")
	now := time.Now()
	lastMonth := now.AddDate(0, -1, 0)
	cardsCmd.Flags().String("updatedSince", lastMonth.Format(time.RFC3339), "filter for cards for the Primary Label")
}
