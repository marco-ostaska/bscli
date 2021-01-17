package squads_test

import (
	"fmt"
	"os"
	"syscall"
	"testing"

	"github.com/marco-ostaska/bscli/cmd"
	"github.com/marco-ostaska/bscli/cmd/vault"
)

func TestSquadsMain(t *testing.T) {
	os.Stdout = os.NewFile(uintptr(syscall.Stdin), "/dev/null")

	t.Run("Testing Command", func(t *testing.T) {

		cmd.RootCmd.SetArgs([]string{"squads"})
		cmd.RootCmd.SilenceErrors = true
		cmd.RootCmd.SilenceUsage = true
		if err := cmd.RootCmd.Execute(); err != nil {
			t.Errorf(err.Error())
		}
	})

	t.Run("Testing Function", func(t *testing.T) {

		if err := squadsTeste(); err != nil {
			t.Errorf(err.Error())
		}
	})
}

func squadsTeste() error {
	vault.ReadVault()

	var gQL struct {
		Data struct {
			Squads []struct {
				ID   string `json:"id"`
				Name string `json:"name"`
			} `json:"squads"`
		} `json:"data"`
	}

	query := `{
		squads {
		  id
		  name
		  squadUsersCount
		}
	  }
	  `

	if err := vault.HTTP.QueryGraphQL(query, &gQL); err != nil {
		if err.Error() == vault.ErrLoginFailure.Error() {
			return fmt.Errorf("Login Failure, please check your token and url and try again")
		}
		return err
	}

	squads := gQL.Data.Squads

	for _, v := range squads {
		if v.Name == "teste graphql" {
			return nil
		}
	}

	return fmt.Errorf("Can not retrieve data")
}
