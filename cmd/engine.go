package commands

import (
	"fmt"
	"os"
	"os/exec"
	"sync"
	"time"
)

type Engine struct {
	KeyChannel  chan string
	DoneChannel chan struct{}
}

func (e Engine) Setup() {
	Log.DebugLog.Printf("Setup Engine \n")
	setupShellSettings()
	e.KeyChannel = make(chan string, 1)
	e.DoneChannel = make(chan struct{})
}

func (e Engine) StartSession() (string, error) {
	Log.DebugLog.Printf("Start Session\n")
	wg := new(sync.WaitGroup)
	wg.Add(3)
	go e.runTimer(10*time.Second, "Session 1", wg)
	go e.inputKeyReceiver(wg)
	go e.inputKeyHandler(wg)
	wg.Wait()
	close(e.KeyChannel)
	close(e.DoneChannel)
	return "Finished", nil
}

func (e Engine) runTimer(duration time.Duration, description string, wg *sync.WaitGroup) {
	Log.DebugLog.Printf("runTimer \n")
	defer wg.Done()
	startTime := time.Now()
	for {
		select {
		case key, ok := <-e.KeyChannel:
			if !ok {
				Log.DebugLog.Printf("keyChannel !ok\n")
				return
			}
			switch key {
			case "b":
				Log.DebugLog.Printf("pressed b\n")
				fmt.Printf("Break")
			case "q":
				Log.DebugLog.Printf("pressed q\n")
				close(e.DoneChannel)
				close(e.KeyChannel)
			}
		default:
			elapsed := time.Since(startTime)
			Log.Stdout.Printf("Elapsed: %s", fmtDuration(elapsed))
			time.Sleep(time.Second)
			if elapsed >= duration {
				Log.Stdout.Println("Timer Stopped")
				return
			}
		}
	}
}

func (e *Engine) inputKeyReceiver(wg *sync.WaitGroup) {
	Log.DebugLog.Printf("inputKeyReceiver\n")
	defer wg.Done()
	var b []byte = make([]byte, 1)
	for {

		select {
		case _, ok := <-e.DoneChannel:
			if !ok {
				Log.DebugLog.Printf("Key Receiver done")
				return
			}
		default:
			Log.Stdout.Printf("Key please:")
			//take care, this will block until key pressed
			os.Stdin.Read(b)
			i := string(b)
			e.KeyChannel <- i
			Log.Stdout.Printf("got key %s", i)
		}
	}
}

func (e *Engine) inputKeyHandler(wg *sync.WaitGroup) {
	Log.DebugLog.Printf("inputKeyHandler\n")
	defer wg.Done()
	for {
		select {

		case key := <-e.KeyChannel:
			Log.DebugLog.Printf("%s pressed", key)
		default:
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
