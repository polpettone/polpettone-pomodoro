package commands

import (
	"fmt"
	"os"
	"os/exec"
	"time"
)

type Engine struct {
	Repo Repo
}

type Session struct {
	Description string        `json:"description"`
	Duration    time.Duration `json:"duration"`
	Start       time.Time     `json:"start"`
}

func (s Session) ToString() string {
	return fmt.Sprintf("%s %s %s", s.Start.Format("2006-01-02 15:04:05"), fmtDuration(s.Duration), s.Description)
}

func (s Session) ToMarkdownTableRow() string {
	return fmt.Sprintf("|%s| %s| %s| %s|\n", s.Start.Format("2006-01-02 15:04:05"), fmtDuration(s.Duration), s.Description, "")
}

func (e *Engine) Setup() {
	setupShellSettings()
	repo := Repo{
		Path: "./sessions.json",
	}
	e.Repo = repo

	Log.DebugLog.Printf("Setup Engine %v", e)
}

func (e *Engine) ExportToMDTable(path string) error {

	sessions, err := e.Repo.Load()
	if err != nil {
		return err
	}

	tableHeader0 := fmt.Sprintf("|Start|Dauer|Beschreibung|Anmerkungen|")
	tableHeader1 := fmt.Sprintf("|-----|-----|------|------|")

	body := ""
	for _, s := range sessions {
		body += fmt.Sprintf("%s", s.ToMarkdownTableRow())
	}

	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return err
	}

	table := fmt.Sprintf("\n\n%s\n%s\n%s\n", tableHeader0, tableHeader1, body)
	f.WriteString(table)

	return nil
}

func (e Engine) StartSession(duration time.Duration, description string, finishCommand string) (string, error) {
	startTime := time.Now()

	session := Session{
		Description: description,
		Duration:    duration,
		Start:       startTime,
	}

	Log.Stdout.Printf("Start Session: %s\n", session.ToString())
	Log.SessionLog.Printf("%s\n", session.ToString())
	e.Repo.Save(session)

	for {
		elapsed := time.Since(startTime)
		Log.Stdout.Printf("session is running: \n %s\n", session.ToString())
		Log.Stdout.Printf("finish command: %s\n", finishCommand)
		Log.Stdout.Printf("elapsed: %s", fmtDuration(elapsed))
		time.Sleep(time.Second)
		clearScreen()
		if elapsed >= duration {
			Log.Stdout.Println("Timer Stopped")
			command := exec.Command(finishCommand)
			command.Run()
			return "finished", nil
		}
	}
}

func fmtDuration(d time.Duration) string {
	hours := int(d.Hours()) % 60
	minute := int(d.Minutes()) % 60
	return fmt.Sprintf("%02d:%02d", hours, minute)
}

func clearScreen() {
	fmt.Print("\033[H\033[2J")
}

func setupShellSettings() {
	// disable input buffering
	exec.Command("stty", "-F", "/dev/tty", "cbreak", "min", "1").Run()
	// do not display entered characters on the screen
	exec.Command("stty", "-F", "/dev/tty", "-echo").Run()
}

func waitForEnter() {
	fmt.Scanln() // Wartet auf Benutzereingabe (Enter)
}
