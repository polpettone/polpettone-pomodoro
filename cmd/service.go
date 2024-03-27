package commands

import (
	"fmt"
	"time"
)

func startSession() (string, error) {
	pomodoroTimer()
	return "started session", nil
}

func pomodoroTimer() {
	runTimer(10*time.Second, "Session 1")
}

func waitForEnter() {
	fmt.Scanln() // Wartet auf Benutzereingabe (Enter)
}

func runTimer(duration time.Duration, description string) {
	startTime := time.Now()
	for {
		elapsed := time.Since(startTime)
		fmt.Printf("Elapsed: %s", fmtDuration(elapsed))
		time.Sleep(time.Second)
		ClearScreen()
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

func manage() {
	shortBreakDuration := 5 * time.Minute
	longBreakDuration := 15 * time.Minute
	pomodoroCount := 0
	for {
		fmt.Printf("Pomodoro %d\n", pomodoroCount+1)
		fmt.Println("Arbeitet für 25 Minuten... (Drücke Enter, um zu pausieren)")
		waitForEnter()
		if pomodoroCount%4 == 3 {
			fmt.Println("Lange Pause für 15 Minuten... (Drücke Enter, um zu pausieren)")
			waitForEnter()
			runTimer(longBreakDuration, "long break")
		} else {
			fmt.Println("Kurze Pause für 5 Minuten... (Drücke Enter, um zu pausieren)")
			waitForEnter()
			runTimer(shortBreakDuration, "short break")
		}
		pomodoroCount++
	}
}

func ClearScreen() {
	fmt.Print("\033[H\033[2J")
}
