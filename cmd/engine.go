package commands

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
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
	return fmt.Sprintf("%s %s %s", s.Start.Format("2006-01-02 15:04:05"), fmtDurationSecondsMinutes(s.Duration), s.Description)
}

func (s Session) ToMarkdownTableRow() string {
	return fmt.Sprintf("|%s| %s| %s| %s|\n", s.Start.Format("2006-01-02 15:04:05"), fmtDurationMinutesHourse(s.Duration), s.Description, "")
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

	Log.InfoLog.Printf("Start Session: %s\n", session.ToString())
	Log.SessionLog.Printf("%s\n", session.ToString())
	e.Repo.Save(session)

	pomodoroStatusFile, err := getPomodoroStatusFile()
	if err != nil {
		return "", err
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for sig := range c {
			Log.Stdout.Printf("%s received. Timer Stopped", sig)
			os.WriteFile(pomodoroStatusFile, []byte("stopped"), 0644)
			os.Exit(0)
		}
	}()

	for {
		elapsed := time.Since(startTime)
		Log.Stdout.Printf("session is running: \n %s\n", session.ToString())
		Log.Stdout.Printf("finish command: %s\n", finishCommand)

		elapsedOutput := fmt.Sprintf("%s: %s", description, fmtDurationSecondsMinutes(elapsed))
		Log.Stdout.Printf(elapsedOutput)

		os.WriteFile(pomodoroStatusFile, []byte(elapsedOutput), 0644)

		time.Sleep(time.Second)
		clearScreen()
		if elapsed >= duration {
			Log.Stdout.Println("Timer Stopped")
			command := exec.Command(finishCommand)
			err := command.Run()
			if err != nil {
				Log.ErrorLog.Print(err)
			}
			os.WriteFile(pomodoroStatusFile, []byte("stopped"), 0644)
			return "finished", nil
		}
	}
}

func getPomodoroStatusFile() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", nil
	}
	return fmt.Sprintf("%s/pomodoro/pomodoro-status", homeDir), nil
}

func fmtDurationMinutesHourse(d time.Duration) string {
	hours := int(d.Hours()) % 60
	minute := int(d.Minutes()) % 60
	return fmt.Sprintf("%02d:%02d", hours, minute)
}

func fmtDurationSecondsMinutes(d time.Duration) string {
	seconds := int(d.Seconds()) % 60
	minute := int(d.Minutes()) % 60
	return fmt.Sprintf("%02d:%02d", minute, seconds)
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
