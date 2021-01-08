package squads_test

import (
	"testing"

	"github.com/marco-ostaska/bscli/cmd"
)

func TestSquads(t *testing.T) {

	cmd.RootCmd.SetArgs([]string{"squads"})
	cmd.RootCmd.SilenceErrors = true
	cmd.RootCmd.SilenceUsage = true
	if err := cmd.RootCmd.Execute(); err != nil {
		t.Errorf(err.Error())
	}
}
