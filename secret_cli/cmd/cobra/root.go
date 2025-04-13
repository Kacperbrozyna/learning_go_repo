package cobra

import (
	"path/filepath"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
)

var encodingKey string

func init() {
	RootCmd.PersistentFlags().StringVarP(&encodingKey, "key", "k", "thisis32bytehardcodedpassphrase!", "the encoding key for encrypting and decrypting")
}

var RootCmd = &cobra.Command{
	Use:   "secret",
	Short: "Secret is an api key and other secrets manager",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func secretsPath() string {
	home, _ := homedir.Dir()
	return filepath.Join(home, "secrets.json")
}
