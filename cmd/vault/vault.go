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

// Package vault is mainly a reference to cobra command vault
//
// But it has the essentials to vault adminstration to be used throughout the application.
package vault

import (
	"fmt"
	"log"
	"strings"

	"github.com/marco-ostaska/httpcalls"
	"github.com/marco-ostaska/uvault"
	"github.com/spf13/cobra"
)

// vault basic constants
const (
	Dir    = "bscli"               // Vault user dir
	File   = "bscli.vlt"           // vault usr file
	APIKey = "Bluesight-API-Token" // default bluesight token key
)

// ErrLoginFailure message when got an login error
var ErrLoginFailure error = fmt.Errorf(`invalid character '<' looking for beginning of value`)

// ReadVault reads the user vault contents
func ReadVault() {
	if err := Credential.ReadFile(Dir, File); err != nil {
		if strings.Contains(err.Error(), "no such file or directory") {
			fmt.Println(newCmd.Long)
			newCmd.Usage()
			fmt.Printf("\n⛔️ ")
			log.Fatalln("No vault created for user, please try to create it using the instruction above first 👀")
		}
		log.Fatalln(err)
	}

	HTTP.URL = Credential.URL
	HTTP.AuthValue = Credential.DecryptedKValue
	HTTP.AuthKey = Credential.APIKey

}

// Credential is a reference to uvault.Credential
var Credential uvault.Credential

// HTTP is a reference to httpcalls.APIData with uservault uploaded
var HTTP httpcalls.APIData

// Cmd represents the vault command
var Cmd = &cobra.Command{
	Use:   "vault",
	Short: "create or update vault credentials",
	Long: `Create or update vault credentials.
For more information how to create a token please 
check the API section https://portal.bluesight.io/tutorial.html 
`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("%v\n", cmd.Long)
		if err := cmd.Usage(); err != nil {
			log.Fatalln(err)
		}
	},
}

func init() {
	Cmd.AddCommand(deleteCmd)

	errs := []error{
		addCommandNewCmd(),
		addCommandUpdateCmd(),
	}

	for _, err := range errs {
		if err != nil {
			log.Fatalln(err)
		}
	}

}
