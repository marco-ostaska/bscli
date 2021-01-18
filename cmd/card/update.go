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

package card

import (
	"fmt"

	"github.com/marco-ostaska/bscli/cmd/vault"
	"github.com/spf13/cobra"
)

type graphQLUpdate struct {
	Data struct {
		UpdateCard struct {
			Card struct {
				Identifier string `json:"identifier"`
				Title      string `json:"title"`
			} `json:"card"`
			Errors []struct {
				Message string `json:"message"`
			} `json:"errors"`
		} `json:"UpdateCard"`
	} `json:"data"`
	Errors []struct {
		Message   string `json:"message"`
		Locations []struct {
			Line   int `json:"line"`
			Column int `json:"column"`
		} `json:"locations"`
		Path []string `json:"path"`
	} `json:"errors"`
}

// cardsCmd represents the cards command
var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "update an existing card",
	Long: `update an existing card
	`,
	Example: `bscli card update --card "09lfdk" -s "Default Swimlane" \
-w "Backlog" \
-t "My new card" \
-d "My new card<br>new line>" \
-p "Primary label1" -p "primary label2" \
-a "assingee@email.com" -a "assingee2@email.com" \
--dueDate "01/31/2021 15:00:00"
`,
	RunE: update,
}

func buildQueryUpdate() (string, error) {

	query := fmt.Sprintln("mutation {")
	query = fmt.Sprintf("%s  updateCard(\n", query)
	query = fmt.Sprintf("%s    input: {\n", query)
	query = fmt.Sprintf(`%s      cardIdentifier: "%s"`+"\n", query, flags.card)
	query = fmt.Sprintf("%s      cardAttributes: {\n", query)
	if len(flags.swimlane) > 0 {
		query = fmt.Sprintf(`%s        swimlaneName: "%s"`+"\n", query, flags.swimlane)
	}

	if len(flags.workstate) > 0 {
		query = fmt.Sprintf(`%s        workstateName: "%s"`+"\n", query, flags.workstate)
	}

	if len(flags.title) > 0 {
		query = fmt.Sprintf(`%s        title: "%s"`+"\n", query, flags.title)
	}

	if len(flags.dueDate) > 0 {
		if err := validateDueDate(flags.dueDate); err != nil {
			return flags.dueDate, err
		}
		query = fmt.Sprintf(`%s        dueDate: "%s"`+"\n", query, flags.dueDate)
	}

	if len(flags.description) > 0 {
		query = fmt.Sprintf(`%s        description: "%s"`+"\n", query, flags.description)
	}

	if len(flags.assignees) > 0 {
		email := fmtStringSlice(flags.assignees)
		query = fmt.Sprintf(`%s        assigneeEmails: %s`+"\n", query, email)
	}

	if len(flags.primarylabel) > 0 {
		primaryLabelNames := fmtStringSlice(flags.primarylabel)
		query = fmt.Sprintf(`%s        primaryLabelNames: %s`+"\n", query, primaryLabelNames)

	}

	endQuery := `     }
   }
  )	
  {
    card {
      identifier
      title
    }
    errors {
      message
    }
  }
}`

	query = fmt.Sprintf("%s %s", query, endQuery)

	return query, nil

}

func update(cmd *cobra.Command, args []string) error {

	query, err := buildQueryUpdate()
	if err != nil {
		return err
	}

	vault.ReadVault()
	var gQL graphQLUpdate

	if err := vault.HTTP.QueryGraphQL(query, &gQL); err != nil {
		return err
	}

	if len(gQL.Errors) > 0 {

		for _, e := range gQL.Errors {
			return fmt.Errorf(e.Message)
		}
	}

	if len(gQL.Data.UpdateCard.Errors) > 0 {

		errCount := 1

		for _, e := range gQL.Data.UpdateCard.Errors {

			switch e.Message {
			case "not found by email":
				return fmt.Errorf("assignee(s) email(s) not found for squad, card not created, please check assignee(s) ans try again")
			case "not found by name":
				return fmt.Errorf("Primary Label(s) not found for squad, card not created, primary label must exist in bluesight")
			case "must exist for the selected workstateName":
				return fmt.Errorf("selected workstateName not found for squad, card not created")
			case "must exist":
				if errCount < len(gQL.Data.UpdateCard.Errors) {
					errCount++
					continue
				}
				return fmt.Errorf("All fields must exist in BlueSight, please check all of them. (Probably wrong Swimlane")
			case "Must be after opened date":
				return fmt.Errorf("dueDate must be higher than opened date")
			default:
				return fmt.Errorf(e.Message)
			}
		}
	}

	fmt.Println("card updated successful")
	fmt.Println("Identifier:", gQL.Data.UpdateCard.Card.Identifier)
	fmt.Println("Title     :", gQL.Data.UpdateCard.Card.Title)

	return nil
}

func init() {

	updateCmd.Flags().StringVar(&flags.card, "card", "", "card identifier")
	updateCmd.Flags().StringVarP(&flags.swimlane, "swimlane", "s", "", "swimlane name")
	updateCmd.Flags().StringVarP(&flags.workstate, "workstate", "w", "", "workstate name")
	updateCmd.Flags().StringVarP(&flags.title, "title", "t", "", "card title")
	updateCmd.Flags().StringVarP(&flags.description, "description", "d", "", "card description")
	updateCmd.Flags().StringSliceVarP(&flags.assignees, "assignees", "a", nil, "card assignee emails")
	updateCmd.Flags().StringSliceVarP(&flags.primarylabel, "primarylabel", "p", nil, "card primary label names")
	updateCmd.Flags().StringVar(&flags.dueDate, "dueDate", "", "card due date")

	if err := updateCmd.MarkFlagRequired("card"); err != nil {
		return
	}

	if err := updateCmd.MarkFlagRequired("swimlane"); err != nil {
		return
	}

}
