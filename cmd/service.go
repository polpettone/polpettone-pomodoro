package commands

import (
	"fmt"
	"os/exec"
	"time"
)

type Engine struct {
}

func (e Engine) Setup() {
	setupShellSettings()
}

func (e Engine) StartSession() (string, error) {
	runTimer(10*time.Second, "Session 1")
	return "started session", nil
}

func runTimer(duration time.Duration, description string) {
	startTime := time.Now()
	for {
		elapsed := time.Since(startTime)
		fmt.Printf("Elapsed: %s", fmtDuration(elapsed))
		time.Sleep(time.Second)
		clearScreen()
		if elapsed >= duration {
			fmt.Println("Timer Stopped")
			break
		}
	}
}

func fmtDuration(d time.Duration) string {
	hour := int(d.Hours())
	minute := int(d.Minutes()) % 60
	second := int(d.Seconds()) % 60
	return fmt.Sprintf("%02d:%02d:%02d\n", hour, minute, second) // 02:09:37
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
