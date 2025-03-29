package cmd

import (
	"fmt"

	database "github.com/Kacperbrozyna/learning_go_repo/cli_task_manager/db"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Lists out all the tasks in your task list",
	Run: func(cmd *cobra.Command, args []string) {

		tasks, err := database.AllTasks()
		if err != nil {
			fmt.Println("Something went wrong:", err.Error())
			return
		}

		if len(tasks) == 0 {
			fmt.Println("You have no tasks to complete! Trying adding some with the add command!")
			return
		}

		fmt.Println("You have the following tasks:")
		for i, task := range tasks {
			fmt.Printf("%d. %s\n", i+1, task.Value)
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
