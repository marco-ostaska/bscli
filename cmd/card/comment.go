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

type graphQLComment struct {
	Data struct {
		CreateComment struct {
			Comment struct {
				CreatedAt string `json:"createdAt"`
				CreatedBy string `json:"createdBy"`
				Body      string `json:"body"`
			} `json:"comment"`
			Card struct {
				Identifier string `json:"identifier"`
			} `json:"card"`
			Errors []struct {
				Message string `json:"message"`
			} `json:"errors"`
		} `json:"createComment"`
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
var commentCmd = &cobra.Command{
	Use:   "comment",
	Short: "add a comment for an existing card",
	Long: `add a comment for an existing card
	`,
	Example: `bscli card comment --card "09lfdk" -s "This is a comment\n this is a new line" 
	`,
	RunE: comment,
}

func buildQueryComment() (string, error) {

	query := fmt.Sprintln("mutation {")
	query = fmt.Sprintf("%s  createComment(\n", query)
	query = fmt.Sprintf("%s    input: {\n", query)
	query = fmt.Sprintf(`%s      cardIdentifier: "%s"`+"\n", query, flags.card)
	query = fmt.Sprintf(`%s      body: "%s"`+"\n", query, flags.comment)

	endQuery := `      }
		)
		{
		comment {
			createdAt
			createdBy
		body
		}
		card {
			identifier
		}
		errors {
			path
			message
		}
		}
		}`

	query = fmt.Sprintf("%s %s", query, endQuery)

	return query, nil

}

func comment(cmd *cobra.Command, args []string) error {

	query, err := buildQueryComment()
	if err != nil {
		return err
	}

	vault.ReadVault()
	var gQL graphQLComment

	if err := vault.HTTP.QueryGraphQL(query, &gQL); err != nil {
		return err
	}

	if len(gQL.Errors) > 0 {

		for _, e := range gQL.Errors {
			return fmt.Errorf(e.Message)
		}
	}

	if len(gQL.Data.CreateComment.Errors) > 0 {

		errCount := 1

		for _, e := range gQL.Data.CreateComment.Errors {

			switch e.Message {
			case "not found by email":
				return fmt.Errorf("assignee(s) email(s) not found for squad, card not created, please check assignee(s) ans try again")
			case "not found by name":
				return fmt.Errorf("Primary Label(s) not found for squad, card not created, primary label must exist in bluesight")
			case "must exist for the selected workstateName":
				return fmt.Errorf("selected workstateName not found for squad, card not created")
			case "must exist":
				if errCount < len(gQL.Data.CreateComment.Errors) {
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

	fmt.Println("card comment created successful")
	fmt.Println("Identifier :", gQL.Data.CreateComment.Card.Identifier)
	fmt.Println("created At :", gQL.Data.CreateComment.Comment.CreatedAt)
	fmt.Println("created By :", gQL.Data.CreateComment.Comment.CreatedBy)
	fmt.Println("Comment    :", gQL.Data.CreateComment.Comment.Body)

	return nil
}

func init() {

	commentCmd.Flags().StringVar(&flags.card, "card", "", "card identifier")
	commentCmd.Flags().StringVarP(&flags.comment, "comment", "c", "", "card comment")

	errs := [2]error{
		commentCmd.MarkFlagRequired("card"),
		commentCmd.MarkFlagRequired("comment"),
	}

	for _, err := range errs {
		if err != nil {
			return
		}
	}

}
