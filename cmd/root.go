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

package cmd

import (
	"fmt"
	"os"

	"github.com/marco-ostaska/bscli/cmd/card"
	"github.com/marco-ostaska/bscli/cmd/squad"
	"github.com/marco-ostaska/bscli/cmd/squads"
	"github.com/marco-ostaska/bscli/cmd/vault"
	"github.com/spf13/cobra"
)

// Version represents bscli version
var version string = "dev"
var versionMsg string = fmt.Sprintf(`bscli %s

Copyright (C) 2021 bscli is released under GNU General Public License v3 
(GPLv3) <http://www.gnu.org/licenses/>

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

Written By Marco Ostaska
`, version)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:     "bscli",
	Short:   "A command line tool for bluesight.io",
	Version: version,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the RootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	// RootCmd.DisableAutoGenTag = true
	// doc.GenMarkdownTree(RootCmd, "./docs")
}

func init() {

	RootCmd.AddCommand(vault.Cmd)
	RootCmd.AddCommand(squads.Cmd)
	RootCmd.AddCommand(squad.Cmd)
	RootCmd.AddCommand(card.Cmd)

	RootCmd.PersistentFlags().BoolP("help", "h", false, "display this help and exit")
	RootCmd.Flags().BoolP("version", "v", false, "output version information and exit")

	RootCmd.SetVersionTemplate(versionMsg)

}
