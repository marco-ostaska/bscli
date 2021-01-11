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

	"github.com/spf13/cobra"
)

type graphQL struct {
	Data struct {
		CreateCard struct {
			Card struct {
				Identifier string `json:"identifier"`
				Title      string `json:"title"`
			} `json:"card"`
			Errors interface{} `json:"errors"`
		} `json:"createCard"`
	} `json:"data"`
}

// cardsCmd represents the cards command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "create a new card for squad",
	Long: `create a new card for squad
	`,
	Example: `bscli card create --squad 01234 -s "Default Swimlane" \
-w "Backlog" \
-t "My new card" 
-d "My new card<br>new line> \
-p "Primary label1" -p "primary label2" \
-a "assingee@email.com" -a "assingee2@email.com" \
--dueDate "01/31/2021 15:00:00"
`,
	RunE: create,
}

func validateDueDate(cmd *cobra.Command) (string, error) {

	layout := "01/02/2006 15:04:05"
	dueDate, err := cmd.Flags().GetString("dueDate")
	if err != nil {
		return dueDate, err
	}
	_, err = time.Parse(layout, dueDate)

	if err != nil {
		return dueDate, fmt.Errorf(`Invalid dueDate, it should be in "MM/dd/yyyy HH:mm:ss" format`)
	}
	return dueDate, err
}

func fmtStringSlice(s []string) string {
	q := "["
	for _, v := range s {
		q = q + fmt.Sprintf(`  
		{ value: "%s" },`, v)
	}

	q = fmt.Sprintf("%s\n        ]", q)

	return q
}

func buildQuery() {

	query := fmt.Sprintln("mutation {")
	query = fmt.Sprintf("%s  createCard(\n", query)
	query = fmt.Sprintf("%s    input: {\n", query)
	query = fmt.Sprintf("%s      cardAttributes: {\n", query)
	id := "012345"
	query = fmt.Sprintf("%s        squadId: %s\n", query, id)
	swimlaneName := "Default swimlaneName"
	query = fmt.Sprintf(`%s        swimlaneName: "%s"`+"\n", query, swimlaneName)
	workstateName := "Backlog"
	query = fmt.Sprintf(`%s        workstateName: "%s"`+"\n", query, workstateName)
	title := "Mais um teste bug email"
	query = fmt.Sprintf(`%s        title: "%s"`+"\n", query, title)
	//dueDate := "01/31/2021 15:00:00"
	dueDate := ""
	if len(dueDate) > 0 {
		query = fmt.Sprintf(`%s        dueDate: "%s"`+"\n", query, dueDate)
	}
	var description string
	//description := "Ola amiguinhos"
	if len(description) > 0 {
		query = fmt.Sprintf(`%s        description: "%s"`+"\n", query, description)
	}

	email := fmtStringSlice([]string{""})

	if len(email) > 0 {
		query = fmt.Sprintf(`%s        assigneeEmails: %s`+"\n", query, email)
	}

	var pr []string

	if pr != nil {
		primaryLabelNames := fmtStringSlice([]string{"L 1", "2 @ email"})
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

	fmt.Println(query)

}

func create(cmd *cobra.Command, args []string) error {

	// dueDate, err := validateDueDate(cmd)

	// if err != nil {
	// 	return err
	// }
	// fmt.Println(dueDate)

	buildQuery()
	// fmt.Println(pLabels)

	// query := `mutation{
	// 	createCard(
	// 	  input: {
	// 		cardAttributes: {
	// 		  squadId: 35106
	// 		  swimlaneName: "Default "
	// 		  workstateName: "Backlog"
	// 		  title: "Mais um teste bug email"
	// 		  dueDate: "01/31/2021 15:00:00"
	// 		  primaryLabelNames: [
	// 			{ value: "bug report"},
	// 			{ value: "bug report2"},
	// 		  ]
	// 		  description: "ola amiguinhos"
	// 		  assigneeEmails: [
	// 			{ value: "marcoan@ccccc"},
	// 			{ value: "jnzaia@cccccc"},
	// 		  ]
	// 		}
	// 	  }
	// 	)
	// 	{
	// 	  card {
	// 		identifier
	// 		title
	// 	  }
	// 	  errors {
	// 		message
	// 	  }
	// 	}
	//   }`

	return nil
}

func init() {

	createCmd.Flags().String("squad", "", "squad id")
	createCmd.Flags().StringP("swimlane", "s", "", "swimlane name")
	createCmd.Flags().StringP("workstate", "w", "", "workstate name")
	createCmd.Flags().StringP("title", "t", "", "card title")
	createCmd.Flags().StringP("description", "d", "", "card description")
	createCmd.Flags().StringSliceP("assignee", "a", nil, "card assignee emails")
	createCmd.Flags().StringSliceP("primarylabel", "p", nil, "card primary label names")
	createCmd.Flags().String("dueDate", "", "card due date")

	// errs := [4]error{
	// 	createCmd.MarkFlagRequired("squad"),
	// 	createCmd.MarkFlagRequired("swimlane"),
	// 	createCmd.MarkFlagRequired("workstate"),
	// 	createCmd.MarkFlagRequired("title"),
	// }

	// for _, err := range errs {
	// 	if err != nil {
	// 		return
	// 	}
	// }

}
