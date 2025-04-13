package cobra

import (
	"MessingWithGo/Learning/learning_go_repo/secret_cli/secret"
	"fmt"

	"github.com/spf13/cobra"
)

var getCmd = &cobra.Command{
	Use:   "Get",
	Short: "Gets a secret from your secret storage",
	Run: func(cmd *cobra.Command, args []string) {
		v := secret.File([]byte(encodingKey), secretsPath())
		key := args[0]
		value, err := v.Get(key)
		if err != nil {
			panic(err)
		}

		fmt.Printf("%s=%s\n", key, value)
	},
}

func init() {
	RootCmd.AddCommand(getCmd)
}
