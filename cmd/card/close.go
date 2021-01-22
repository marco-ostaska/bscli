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
	"time"

	"github.com/marco-ostaska/bscli/cmd/vault"
	"github.com/spf13/cobra"
)

type graphQLClose struct {
	Data struct {
		Squad struct {
			ID    string `json:"id"`
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
		UpdateCard struct {
			Card struct {
				Identifier    string `json:"identifier"`
				Title         string `json:"title"`
				WorkstateType string `json:"workstateType"`
				Swimlane      string `json:"swimlane"`
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

// closeCMD represents the cards command
var closeCmd = &cobra.Command{
	Use:   "close",
	Short: "close an existing card",
	Long: `close an existing card
	`,
	Example: `bscli card close --card [card identifier] --squad [squad id]
	`,
	RunE: closeCard,
}

func closeCard(cmd *cobra.Command, args []string) error {

	id, err := cmd.Flags().GetString("squad")
	if err != nil {
		return err
	}

	cid, err := cmd.Flags().GetString("card")
	if err != nil {
		return err
	}

	vault.ReadVault()

	now := time.Now()
	lastMonth := now.AddDate(0, -1, 0)

	query := fmt.Sprintf(`{
		squad(id: %s){
			cards(updatedSince: "%v", closed: false){
			identifier
			title
			swimlane
			workstateType
		    }   
		}
	}`, id, lastMonth.Format(time.RFC3339))

	var gQL graphQLClose

	if err := vault.HTTP.QueryGraphQL(query, &gQL); err != nil {
		return err
	}

	if err := gQL.close(cid); err != nil {
		cmd.SilenceErrors = true
		cmd.SilenceUsage = true
		return err
	}

	return nil
}

func (gQL graphQLClose) close(id string) error {

	query, err := gQL.buildQueryClose(id)
	if err != nil {
		return err
	}

	if err := vault.HTTP.QueryGraphQL(query, &gQL); err != nil {
		return err
	}

	fmt.Println("card closed successful")

	return nil
}

func (gQL graphQLClose) buildQueryClose(id string) (string, error) {
	var query string

	for _, c := range gQL.Data.Squad.Cards {
		if c.Identifier == id {
			query := fmt.Sprintln("mutation {")
			query = fmt.Sprintf("%s  updateCard(\n", query)
			query = fmt.Sprintf("%s    input: {\n", query)
			query = fmt.Sprintf(`%s      cardIdentifier: "%s"`+"\n", query, id)
			query = fmt.Sprintf("%s      cardAttributes: {\n", query)
			query = fmt.Sprintf(`%s        swimlaneName: "%s"`+"\n", query, c.Swimlane)
			query = fmt.Sprintf(`%s        workstateName: "Closed"`+"\n", query)

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
	}

	return query, fmt.Errorf("Card not found or already closed")

}

func init() {

	Cmd.AddCommand(closeCmd)
	closeCmd.Flags().String("card", "", "card identifier")
	closeCmd.Flags().String("squad", "", "squad id")

	if err := closeCmd.MarkFlagRequired("card"); err != nil {
		return
	}

	if err := closeCmd.MarkFlagRequired("squad"); err != nil {
		return
	}

}
