package commands

import (
	"fmt"
	"os/exec"
	"sync"
	"time"
)

type Engine struct {
	KeyChannel chan string
}

func (e Engine) Setup() {
	setupShellSettings()
	e.KeyChannel = make(chan string, 1)
}

func (e Engine) StartSession() (string, error) {
	wg := new(sync.WaitGroup)
	wg.Add(1)
	go runTimer(3*time.Second, "Session 1", wg)
	wg.Wait()
	return "Finished", nil
}

func runTimer(duration time.Duration, description string, wg *sync.WaitGroup) {
	defer wg.Done()
	startTime := time.Now()
	for {
		select {
		default:
			elapsed := time.Since(startTime)
			fmt.Printf("Elapsed: %s", fmtDuration(elapsed))
			time.Sleep(time.Second)
			clearScreen()
			if elapsed >= duration {
				fmt.Println("Timer Stopped")
				return
			}
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
