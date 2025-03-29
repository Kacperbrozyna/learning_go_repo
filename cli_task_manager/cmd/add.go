package cmd

import (
	"fmt"
	"strings"

	database "github.com/Kacperbrozyna/learning_go_repo/cli_task_manager/db"
	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Adds a task to your task list",
	Run: func(comd *cobra.Command, args []string) {
		task := strings.Join(args, " ")

		_, err := database.CreateTask(task)
		if err != nil {
			fmt.Println("Something went wrong:", err)
			return
		}

		fmt.Printf("Added \"%s\" to your task list. \n", task)
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
}
