package squad_test

import (
	"os"
	"testing"

	"github.com/marco-ostaska/bscli/cmd"
)

func TestDisplayUsers(t *testing.T) {
	id := os.Getenv("BS_TEST_ID")

	cmd.RootCmd.SetArgs([]string{"squad", "--id", id, "users"})
	cmd.RootCmd.SilenceErrors = true
	cmd.RootCmd.SilenceUsage = true
	if err := cmd.RootCmd.Execute(); err != nil {
		t.Errorf(err.Error())
	}

}

func TestDisplayswimlane(t *testing.T) {
	id := os.Getenv("BS_TEST_ID")
	tt := []struct {
		name string
		args []string
	}{
		{"calling alias sl", []string{"squad", "--id", id, "sl"}},
		{"calling swimlane", []string{"squad", "--id", id, "swimlane"}},
		{"Testing filter", []string{"squad", "--id", id, "swimlane", "--filterSL", "Default SwimLane"}},
		{"Testing tree", []string{"squad", "--id", id, "swimlane", "--tree"}},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			cmd.RootCmd.SetArgs(tc.args)
			cmd.RootCmd.SilenceErrors = true
			cmd.RootCmd.SilenceUsage = true
			if err := cmd.RootCmd.Execute(); err != nil {
				t.Errorf(err.Error())
			}
		})
	}

}
