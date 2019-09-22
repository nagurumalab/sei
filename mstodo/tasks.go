package mstodo

import (
	"encoding/json"
	"fmt"
	"net/http"
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
type Tasks struct {
	Value    []Task `json:"value"`
	NextLink string `json:"@odata.nextLink"`
}

//ListTasks lists all tasks
func ListTasks(client *http.Client, folderID string, hideCompleted bool) *Tasks {
	t := new(Tasks)
	var resp *http.Response
	var err error
	filter := ""
	if hideCompleted {
		filter += " status eq 'notStarted' "
	}
	params := map[string]string{}
	if filter != "" {
		params["$filter"] = filter
	}
	if folderID != "" {
		reqURL := constructURL([]string{"taskFolders", folderID, "tasks"}, params)
		//log.Printf("Requesting url - %s", reqURL)
		resp, err = client.Get(reqURL)
		if err != nil {
			panic(err)
		}
	}
	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(t)
	if err != nil {
		panic(err)
	}
	return t
}

//Print prints the folders struct to console
func (t Tasks) Print(detailed bool) {
	var err error
	for i, task := range t.Value {
		//TODO: Work on better detailed format
		if detailed {
			_, err = fmt.Printf("%d - %s\n\tID: %s\n\tBody: %s\n",
				i+1, task.Subject, task.ID, task.Body)
		} else {
			_, err = fmt.Printf("%d - %s\n", i+1, task.Subject)
		}

		if err != nil {
			panic(err)
		}
	}
}
