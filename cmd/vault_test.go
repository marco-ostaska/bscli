package cmd

import (
	"testing"
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
			rootCmd.SetArgs(tc.args)
			rootCmd.SilenceErrors = true
			rootCmd.SilenceUsage = true
			if err := rootCmd.Execute(); err != nil {
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
		if err := vCredential.ReadFile(vaultDir, vaultFile); err != nil {
			t.Errorf(err.Error())
		}

		switch {
		case vCredential.APIKey != apiKey:
			t.Errorf("got %v, expected %v", vCredential.APIKey, apiKey)
		case vCredential.DecryptedKValue != keyValue:
			t.Errorf("got %v, expected %v", vCredential.DecryptedKValue, keyValue)
		case vCredential.URL != uri:
			t.Errorf("got %v, expected %v", vCredential.URL, uri)
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
			rootCmd.SetArgs(tc.args)
			rootCmd.SilenceErrors = true
			rootCmd.SilenceUsage = true
			if err := rootCmd.Execute(); err != nil {
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
		if err := vCredential.ReadFile(vaultDir, vaultFile); err != nil {
			t.Errorf(err.Error())
		}

		switch {
		case vCredential.APIKey != apiKey:
			t.Errorf("got %v, expected %v", vCredential.APIKey, apiKey)
		case vCredential.DecryptedKValue != keyValue:
			t.Errorf("got %v, expected %v", vCredential.DecryptedKValue, keyValue)
		case vCredential.URL != uri:
			t.Errorf("got %v, expected %v", vCredential.URL, uri)
		}

	})
}

func TestDeleteVault(t *testing.T) {

	rootCmd.SetArgs([]string{"vault", "delete"})
	if err := rootCmd.Execute(); err != nil {
		t.Errorf(err.Error())
	}

}
