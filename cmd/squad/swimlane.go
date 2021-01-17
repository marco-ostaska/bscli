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

// swimlaneCmd represents the squad command
var swimlaneCmd = &cobra.Command{
	Use:           "swimlane",
	Aliases:       []string{"sl"},
	Short:         "display the swimlane workstates of given squad",
	SilenceErrors: true,
	Example:       `bscli squad --id [squad id] swimlane --filterSL "Default Swimlane" `,
	Long: `display the users for the squad
	`,
	RunE: displayswimlane,
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

func displayswimlane(cmd *cobra.Command, args []string) error {
	vault.ReadVault()
	id, err := cmd.Flags().GetString("id")
	if err != nil {
		return err
	}

	var gQL graphQL
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

	var SwimlaneWorkstatesTotal int
	var BacklogWorkstatesTotal int
	var ActiveWorkstatesTotal int
	var WaitWorkstatesTotal int

	filterSL, err := cmd.Flags().GetString("filterSL")
	if err != nil {
		return err
	}
	tree, err := cmd.Flags().GetBool("tree")
	if err != nil {
		return err
	}

	for _, sl := range gQL.Data.Squad.SwimlaneWorkstates {

		if len(filterSL) > 0 && filterSL != sl.Name {
			continue
		}

		fmt.Println(sl.Name)

		if tree {
			printTree("Backlog WorkStates", sl.BacklogWorkstates)
			printTree("Active WorkStates", sl.ActiveWorkstates)
			printTree("Wait WorkStates", sl.WaitWorkstates)
			fmt.Println()
			SwimlaneWorkstatesTotal++
			BacklogWorkstatesTotal = len(sl.BacklogWorkstates)
			ActiveWorkstatesTotal = len(sl.ActiveWorkstates)
			WaitWorkstatesTotal = len(sl.WaitWorkstates)
		}

	}

	if tree {
		fmt.Printf("%v SwimlaneWorkstates, %v BacklogWorkstates, %v ActiveWorkstates,  %v WaitWorkstates\n",
			SwimlaneWorkstatesTotal, BacklogWorkstatesTotal, ActiveWorkstatesTotal, WaitWorkstatesTotal)
	}

	return nil
}

func init() {
	Cmd.AddCommand(swimlaneCmd)
	swimlaneCmd.Flags().String("filterSL", "", "filter SwimLane")
	swimlaneCmd.Flags().Bool("tree", false, "SwimLane tree with Workstates details")
}
