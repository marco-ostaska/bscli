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

// Package mood defines the mood marbles command to bew called by cobra
package mood

import (
	"fmt"

	"github.com/marco-ostaska/bscli/cmd/vault"
	"github.com/spf13/cobra"
)

type graphQL struct {
	Data struct {
		Squad struct {
			Mood struct {
				MoodRecord []struct {
					CreatedAt string  `json:"createdAt"`
					MoodScore float64 `json:"moodScore"`
					Mood      string  `json:"mood"`
					Comment   string  `json:"comment"`
				} `json:"moodRecord"`
				MoodReport struct {
					Last24Hours     float64 `json:"last24Hours"`
					Last7Days       float64 `json:"last7Days"`
					LastMonth       float64 `json:"lastMonth"`
					MonthlyAverages []struct {
						Date  string  `json:"date"`
						Value float64 `json:"value"`
					} `json:"monthlyAverages"`
				} `json:"moodReport"`
			} `json:"mood"`
		} `json:"squad"`
	} `json:"data"`
}

var flags = []struct {
	Name        string
	Description string
}{
	{"record", "dsplay mood marbles records for given squad"},
	{"report", "dsplay mood marbles report for given squad"},
	{"comments", "display last mood marbles comments for given squad"},
}

// Cmd represents the mood.go command
var Cmd = &cobra.Command{
	Use:   "mood [id]...",
	Args:  cobra.ExactArgs(1),
	Short: "display mood marbles information for a given squad",
	Long: `display mood marbles information for a given squad
	`,
	RunE: moodMain,
}

func moodMain(cmd *cobra.Command, args []string) error {

	if cmd.Flags().NFlag() < 1 {
		return fmt.Errorf("No flag(s) received")
	}
	vault.ReadVault()

	for _, f := range flags {
		if cmd.Flag(f.Name).Changed {
			var gQL graphQL
			switch f.Name {
			case "record":
				if err := gQL.displayRecords(args[0]); err != nil {
					return err
				}
			case "comments":
				if err := gQL.displayComments(args[0]); err != nil {
					return err
				}
			case "report":
				if err := gQL.displayReport(args[0]); err != nil {
					return err
				}
			}

		}
	}
	return nil

}

func (gQL graphQL) displayRecords(id string) error {
	query := fmt.Sprintf(`{
		squad(id: %s){
		  mood{
			moodRecord{
			  createdAt
			  moodScore
			  mood
			  comment
			}
		  }
		}
	  }`, id)

	if err := vault.HTTP.QueryGraphQL(query, &gQL); err != nil {
		return err
	}

	for _, m := range gQL.Data.Squad.Mood.MoodRecord {
		fmt.Println("\ncreated at   :", m.CreatedAt)
		fmt.Println("Score        :", m.MoodScore)
		fmt.Println("Mood         :", m.Mood)
		fmt.Println("Emoji        :", emoji(m.MoodScore))
		if len(m.Comment) > 0 {
			fmt.Println("Comment      :", m.Comment)
		}

	}

	return nil
}

func (gQL graphQL) displayComments(id string) error {
	query := fmt.Sprintf(`{
		squad(id: %s){
		  mood{
			moodRecord{
			  comment
			}
		  }
		}
	  }`, id)

	if err := vault.HTTP.QueryGraphQL(query, &gQL); err != nil {
		return err
	}

	for _, m := range gQL.Data.Squad.Mood.MoodRecord {
		if len(m.Comment) > 0 {
			fmt.Println("-", m.Comment)
		}

	}

	return nil
}

func (gQL graphQL) displayReport(id string) error {
	query := fmt.Sprintf(`{
		squad(id: %s){
		  mood{
			moodReport{
			  last24Hours
			  last7Days
			  lastMonth
			  monthlyAverages{
				date
				value
			  }
			}
		  }
		}
	  }`, id)

	if err := vault.HTTP.QueryGraphQL(query, &gQL); err != nil {
		return err
	}

	mReport := gQL.Data.Squad.Mood.MoodReport

	fmt.Printf("Last 24 hours: %v %s\n", mReport.Last24Hours, emoji(mReport.Last24Hours))
	fmt.Printf("Last 7 Days  : %v %s\n", mReport.Last7Days, emoji(mReport.Last7Days))
	fmt.Printf("Last Month   : %v %s\n", mReport.LastMonth, emoji(mReport.LastMonth))

	for _, m := range mReport.MonthlyAverages {
		fmt.Printf("%-25s: %v %s\n", m.Date, m.Value, emoji(m.Value))
	}

	// for _, m := range gQL.Data.Squad.Mood.MoodReport.MonthlyAverages {
	// 	return

	// }

	return nil
}

func emoji(score float64) string {
	switch {
	case score >= 0 && score < 2:
		return "😩"
	case score >= 2 && score < 3:
		return "😖"
	case score >= 3 && score < 4:
		return "😐"
	case score >= 4 && score < 5:
		return "😆"
	case score >= 5:
		return "😁"
	}

	return ""

}

func init() {

	for _, f := range flags {
		Cmd.Flags().Bool(f.Name, false, f.Description)
	}

}
