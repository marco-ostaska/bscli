package vault_test

import (
	"testing"

	"github.com/marco-ostaska/bscli/cmd"
	"github.com/marco-ostaska/bscli/cmd/vault"
)

func TestNewVault(t *testing.T) {
	tt := []struct {
		name     string
		args     []string
		expected string
	}{
		{"Missing args 1", []string{"vault", "new"}, `required flag(s) "key", "url" not set`},
		{"Missing args 2", []string{"vault", "new", "-k", "!@#$%^&*key"}, `required flag(s) "url" not set`},
		{"Ok", []string{"vault", "new", "-k", "!@#$%^&*key", "--url", "https://xyz.io"}, ""},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {

			cmd.RootCmd.SetArgs(tc.args)
			cmd.RootCmd.SilenceErrors = true
			cmd.RootCmd.SilenceUsage = true
			if err := cmd.RootCmd.Execute(); err != nil {
				if err.Error() == tc.expected {
					t.Skip(err)
					return
				}
				t.Errorf("got: %s, expected: %s", err.Error(), tc.expected)
			}
		})
	}

	t.Run("Check credentials", func(t *testing.T) {
		keyValue := "!@#$%^&*key"
		uri := "https://xyz.io"
		if err := vault.Credential.ReadFile(vault.Dir, vault.File); err != nil {
			t.Errorf(err.Error())
		}

		switch {
		case vault.Credential.APIKey != vault.APIKey:
			t.Errorf("got %v, expected %v", vault.Credential.APIKey, vault.APIKey)
		case vault.Credential.DecryptedKValue != keyValue:
			t.Errorf("got %v, expected %v", vault.Credential.DecryptedKValue, keyValue)
		case vault.Credential.URL != uri:
			t.Errorf("got %v, expected %v", vault.Credential.URL, uri)
		}

	})

}

func TestUpdateVault(t *testing.T) {
	tt := []struct {
		name     string
		args     []string
		expected string
	}{
		{"Missing args 1", []string{"vault", "update"}, `required flag(s) "key" not set`},
		{"OK", []string{"vault", "update", "-k", "mykey"}, `required flag(s) "url" not set`},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			cmd.RootCmd.SetArgs(tc.args)
			cmd.RootCmd.SilenceErrors = true
			cmd.RootCmd.SilenceUsage = true
			if err := cmd.RootCmd.Execute(); err != nil {
				if err.Error() == tc.expected {
					t.Skip(err)
					return
				}
				t.Errorf("got: %s, expected: %s", err.Error(), tc.expected)
			}
		})
	}

	t.Run("Check credentials", func(t *testing.T) {
		keyValue := "mykey"
		uri := "https://xyz.io"
		if err := vault.Credential.ReadFile(vault.Dir, vault.File); err != nil {
			t.Errorf(err.Error())
		}

		switch {
		case vault.Credential.APIKey != vault.APIKey:
			t.Errorf("got %v, expected %v", vault.Credential.APIKey, vault.APIKey)
		case vault.Credential.DecryptedKValue != keyValue:
			t.Errorf("got %v, expected %v", vault.Credential.DecryptedKValue, keyValue)
		case vault.Credential.URL != uri:
			t.Errorf("got %v, expected %v", vault.Credential.URL, uri)
		}

	})
}

func TestDeleteVault(t *testing.T) {

	cmd.RootCmd.SetArgs([]string{"vault", "delete"})
	if err := cmd.RootCmd.Execute(); err != nil {
		t.Errorf(err.Error())
	}

}
