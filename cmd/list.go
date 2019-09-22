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
	"github.com/nagurumalab/learning-go/pani/mstodo"

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
		var folder *mstodo.Folder
		client := mstodo.NewClient()
		folderName, _ := cmd.Flags().GetString("folder-name")
		folderID, _ := cmd.Flags().GetString("folder-id")
		folders := mstodo.ListFolders(client)
		if folderID != "" {
			folder = folders.GetFolderFromID(folderID)
		}
		if folder == nil && folderName != "" {
			folder = folders.GetFolderFromName(folderName)
		}
		if folder == nil {
			folder = folders.GetDefaultFolder()
		}
		//log.Printf("Getting task of folder - %s", folder.ID)
		folder.GetTasks(client).Print(false)
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
		detailed, _ := cmd.Flags().GetBool("detailed")
		mstodo.ListFolders(mstodo.NewClient()).Print(detailed)
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
	listFolders.Flags().BoolP("detailed", "d", false, "Show details of the folders")
	listTasks.Flags().StringP("folder-name", "n", "", "Name of the folder to list the tasks")
	listTasks.Flags().StringP("folder-id", "i", "", "ID of the folder to list the tasks")
}
