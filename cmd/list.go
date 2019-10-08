/*
Copyright © 2019 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"github.com/nagurumalab/sei/mstodo"

	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listTasks = &cobra.Command{
	Use:   "lt",
	Short: "Lists tasks",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		client := mstodo.NewClient()
		var folder *mstodo.Folder
		// TODO: Handle err from all the following flags.Get
		if defaultFolder, _ := cmd.Flags().GetBool("default-folder"); defaultFolder {
			folder = mstodo.ListFolders(client).GetDefaultFolder()
		} else if folderID, _ := cmd.Flags().GetString("folder-id"); folderID != "" {
			folder = mstodo.ListFolders(client).GetFolderFromID(folderID)
		} else if folderName, _ := cmd.Flags().GetString("folder-name"); folder == nil && folderName != "" {
			folder = mstodo.ListFolders(client).GetFolderFromName(folderName)
		}

		//log.Printf("Getting task of folder - %s", folder.ID)
		var tasks []mstodo.Task
		limit, _ := cmd.Flags().GetInt("limit")
		expandedView, _ := cmd.Flags().GetBool("expanded-view")
		hideCompleted := true
		if folder != nil {
			tasks = folder.GetTasks(client, hideCompleted, limit)
		} else {
			tasks = mstodo.ListTasks(client, "", hideCompleted, limit)
		}
		mstodo.PrintTasks(tasks, expandedView)

	},
}

var listFolders = &cobra.Command{
	Use:   "lf",
	Short: "Lists folders",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Move the newclient to a general place. may be pre run or init hook
		expandedView, _ := cmd.Flags().GetBool("expanded-view")
		mstodo.ListFolders(mstodo.NewClient()).Print(expandedView)
	},
}

func init() {
	rootCmd.AddCommand(listTasks, listFolders)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	listFolders.Flags().BoolP("expanded-view", "e", false, "Show details of each folders")

	// TODO: Either default takes precidence over id , which takes over name. Enforce this in the flags
	listTasks.Flags().StringP("folder-name", "n", "", "Name of the folder to list the tasks")
	listTasks.Flags().StringP("folder-id", "i", "", "ID of the folder to list the tasks")
	listTasks.Flags().BoolP("default-folder", "d", false, "List tasks of default folder")
	listTasks.Flags().IntP("limit", "l", 20, "Number of tasks to show (20 default)")
	listTasks.Flags().BoolP("expanded-view", "e", false, "Show details of each tasks")
}
