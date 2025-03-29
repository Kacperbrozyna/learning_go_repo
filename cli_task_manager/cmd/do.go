package cmd

import (
	"fmt"
	"strconv"

	database "github.com/Kacperbrozyna/learning_go_repo/cli_task_manager/db"
	"github.com/spf13/cobra"
)

var doCmd = &cobra.Command{
	Use:   "do",
	Short: "Complete a task in your task list based on the index provided",
	Run: func(cmd *cobra.Command, args []string) {
		var ids []int
		for _, arg := range args {
			id, err := strconv.Atoi(arg)
			if err != nil {
				fmt.Println("Failed to parse the argument: ", arg)
				continue
			}

			ids = append(ids, id)
		}

		tasks, err := database.AllTasks()
		if err != nil {
			fmt.Println("Something went wrong:", err)
			return
		}

		for _, id := range ids {
			if id <= 0 || id > len(tasks) {
				fmt.Println("Invalid task number:", id)
				continue
			}
			task := tasks[id-1]
			err := database.DeleteTask(task.Key)
			if err != nil {
				fmt.Printf("Failed to mark \"%d\" as completed, Error: %s\n", id, err.Error())
				return
			}

			fmt.Printf("Marked \"%d\" as completed. \n", id)
		}
	},
}

func init() {
	rootCmd.AddCommand(doCmd)
}
