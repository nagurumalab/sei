package mstodo

import (
	"encoding/json"
	"fmt"
	"net/http"
)

//Folder Task folder
type Folder struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	IsDefault bool   `json:"isDefaultFolder"`
}

//Folders Task folders from api
type Folders struct {
	Value    []Folder `json:"value"`
	NextLink string   `json:"@odata.nextLink"`
}

//ListFolders lists all the folder for the user
func ListFolders(client *http.Client) *Folders {
	// log.Println("Calling url - ", URLS["ListFolders"].url)
	resp, err := client.Get("https://graph.microsoft.com/beta/me/outlook/taskFolders")
	if err != nil {
		panic(err)
	}
	// log.Println("Response Status - ", resp.Status)
	defer resp.Body.Close()
	f := new(Folders)
	err = json.NewDecoder(resp.Body).Decode(f)
	if err != nil {
		panic(err)
	}
	return f
}

//Print prints the folders struct to console
func (f Folders) Print(detailed bool) {
	var err error
	for i, folder := range f.Value {
		defaultFolder := ""
		if folder.IsDefault {
			defaultFolder = "(Default)"
		}

		if detailed {
			_, err = fmt.Printf("%d - %s %s\n\tID: %s\n\n",
				i+1, folder.Name, defaultFolder, folder.ID)
		} else {
			_, err = fmt.Printf("%d - %s %s\n", i+1, folder.Name, defaultFolder)
		}

		if err != nil {
			panic(err)
		}
	}
}

//GetDefaultFolder loops over and get the default folders
//TODO: Cache the default folder id
func (f Folders) GetDefaultFolder() *Folder {
	for _, folder := range f.Value {
		if folder.IsDefault {
			return &folder
		}
	}
	return nil
}

//GetFolderFromID gets you the folder given an ID
func (f Folders) GetFolderFromID(id string) *Folder {
	if id != "" {
		for _, folder := range f.Value {
			if folder.ID == id {
				return &folder
			}
		}
	}
	return nil
}

//GetFolderFromName gets you the folder given an name pattern
//TODO: Need to support regex in patterns
func (f Folders) GetFolderFromName(namePattern string) *Folder {
	if namePattern != "" {
		for _, folder := range f.Value {
			if folder.Name == namePattern {
				return &folder
			}
		}
	}
	return nil
}

//GetTasks gets tasks under Folder f
func (f Folder) GetTasks(client *http.Client) *Tasks {
	return ListTasks(client, f.ID, true)
}
