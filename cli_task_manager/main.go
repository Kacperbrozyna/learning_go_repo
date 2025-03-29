package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/Kacperbrozyna/learning_go_repo/cli_task_manager/cmd"
	database "github.com/Kacperbrozyna/learning_go_repo/cli_task_manager/db"
	"github.com/mitchellh/go-homedir"
)

func main() {
	home_dir_path, err := homedir.Dir()
	if err != nil {
		onError(err)
	}

	db_path := filepath.Join(home_dir_path, "tasks.db")
	err = database.Init(db_path)
	if err != nil {
		onError(err)
	}

	cmd.Execute()
}

func onError(err error) {
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
