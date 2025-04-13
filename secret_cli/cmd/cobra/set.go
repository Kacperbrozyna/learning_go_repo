package cobra

import (
	"MessingWithGo/Learning/learning_go_repo/secret_cli/secret"
	"fmt"

	"github.com/spf13/cobra"
)

var setCmd = &cobra.Command{
	Use:   "Set",
	Short: "Sets a secret in your secret storage",
	Run: func(cmd *cobra.Command, args []string) {
		v := secret.File([]byte(encodingKey), secretsPath())
		key, value := args[0], args[1]

		err := v.Set(key, value)
		if err != nil {
			panic(err)
		}

		fmt.Printf("Value Set\n")
	},
}

func init() {
	RootCmd.AddCommand(setCmd)
}
