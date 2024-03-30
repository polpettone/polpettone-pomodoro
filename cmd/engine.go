package commands

import (
	"fmt"
	"os/exec"
	"time"
)

type Engine struct {
	KeyChannel  chan string
	DoneChannel chan struct{}
}

type Session struct {
	Description string
	Duration    time.Duration
	Start       time.Time
}

func (s Session) ToString() string {
	return fmt.Sprintf("%s %s %s", s.Start.Format("2006-01-02 15:04:05"), fmtDuration(s.Duration), s.Description)
}

func (e Engine) Setup() {
	Log.DebugLog.Printf("Setup Engine \n")
	setupShellSettings()
	e.KeyChannel = make(chan string, 1)
	e.DoneChannel = make(chan struct{})
}

func (e Engine) StartSession(duration time.Duration, description string) (string, error) {
	startTime := time.Now()

	session := Session{
		Description: description,
		Duration:    duration,
		Start:       startTime,
	}

	Log.Stdout.Printf("Start Session: %s\n", session.ToString())
	Log.SessionLog.Printf("%s\n", session.ToString())

	for {
		elapsed := time.Since(startTime)
		Log.Stdout.Printf("session is running: \n %s\n", session.ToString())
		Log.Stdout.Printf("elapsed: %s", fmtDuration(elapsed))
		time.Sleep(time.Second)
		clearScreen()
		if elapsed >= duration {
			Log.Stdout.Println("Timer Stopped")
			command := exec.Command("i3lock")
			command.Run()
			return "finished", nil
		}
	}
}

func fmtDuration(d time.Duration) string {
	minute := int(d.Minutes()) % 60
	second := int(d.Seconds()) % 60
	return fmt.Sprintf("%02d:%02d", minute, second) // 02:09:37
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
