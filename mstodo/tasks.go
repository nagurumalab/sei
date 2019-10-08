package mstodo

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

//Task struct
type Task struct {
	ID      string `json:"id"`
	Subject string `json:"subject"`
	Body    struct {
		ContentType string `json:"contentType"`
		Content     string `json:"content"`
	} `json:"body"`
	CreatedDateTime string `json:"createdDateTime"`
	DueDateTime     struct {
		DateTime string `json:"dateTime"`
		TimeZone string `json:"timeZone"`
	} `json:"dueDateTime"`
	ReminderDateTime struct {
		DateTime string `json:"dateTime"`
		TimeZone string `json:"timeZone"`
	} `json:"reminderDateTime"`
	IsReminderOn   bool   `json:"isReminderOn"`
	ParentFolderID string `json:"parentFolderId"`
	Status         string `json:"status"`
}

//Tasks list of tasks
type _task struct {
	Value    []Task `json:"value"`
	NextLink string `json:"@odata.nextLink"`
}

//ListTasks lists all tasks
func ListTasks(client *http.Client, folderID string, hideCompleted bool, limit int) []Task {
	t := new(_task)
	var reqURL string
	var resp *http.Response
	var err error
	params := map[string]string{}

	filter := ""
	if hideCompleted {
		filter += " status eq 'notStarted' "
	}
	if filter != "" {
		params["$filter"] = filter
	}
	params["$top"] = strconv.Itoa(limit)

	if folderID != "" {
		reqURL = constructURL([]string{"taskFolders", folderID, "tasks"}, params)
	} else {
		reqURL = constructURL([]string{"tasks"}, params)
	}
	resp, err = client.Get(reqURL)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(t)
	if err != nil {
		panic(err)
	}
	return t.Value
}

//PrintTasks prints the folders struct to console
func PrintTasks(tasks []Task, detailed bool) {
	var err error
	for i, task := range tasks {
		//TODO: Work on better detailed format
		if detailed {
			_, err = fmt.Printf("%d - %s\n\tID: %s\n\tFolder: %s\n\tCreated: %s\n\tBody: %s\n\n",
				i+1, task.Subject, task.ID, task.ParentFolderID, task.CreatedDateTime, task.Body.Content)
		} else {
			_, err = fmt.Printf("%d - %s\n", i+1, task.Subject)
		}

		if err != nil {
			panic(err)
		}
	}
}
