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

// Cmd represents the squad command
var Cmd = &cobra.Command{
	Use:           "squad",
	Short:         "display information for a given squad",
	SilenceErrors: true,
	Long: `display information for a given squad
	`,
}

func init() {

	Cmd.PersistentFlags().String("id", "", "squad id")
	//vault.ReadVault()
}
