package commands

import (
	"encoding/json"
	"os"
)

type Repo struct {
	Path string
}

func (repo Repo) Save(session Session) error {

	var sessions []Session
	sessions, err := repo.Load()

	if err != nil {
		Log.InfoLog.Printf("Could not Load sessions.json, new one will be created  %s", err)
		sessions = []Session{}
	}

	Log.DebugLog.Printf("Loaded Sessions %v", sessions)

	sessions = append(sessions, session)

	j, err := json.Marshal(sessions)
	if err != nil {
		return err
	}

	f, err := openFile(repo.Path)
	if err != nil {
		return err
	}

	_, err = f.Write(j)
	if err != nil {
		return err
	}

	return nil
}

func (repo Repo) Load() ([]Session, error) {

	Log.DebugLog.Printf("Load Sessions from Path: %s", repo.Path)

	j, err := os.ReadFile(repo.Path)

	var sessions []Session
	json.Unmarshal(j, &sessions)

	if err != nil {
		return nil, err
	}

	return sessions, nil
}

func openFile(path string) (*os.File, error) {
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		Log.ErrorLog.Printf("error opening file: %s, %v", path, err)
		return nil, err
	}
	return f, nil
}
