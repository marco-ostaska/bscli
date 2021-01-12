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

	"github.com/spf13/cobra"
)

// Cmd represents the cards command
var Cmd = &cobra.Command{
	Use:   "card",
	Short: "create, update or create comment for a given card",
	Long: `create, update or create comment for a given card
	`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(cmd.Long)
		cmd.Usage()
	},
}

func init() {
	Cmd.AddCommand(createCmd)
	Cmd.AddCommand(updateCmd)
	Cmd.AddCommand(commentCmd)
}
