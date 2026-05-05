package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
)

const (
	Reset  = "\033[0m"
	Red    = "\033[31m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
	Blue   = "\033[34m"
	Cyan   = "\033[36m"
	Gray   = "\033[37m"
)

const banner = `
__/\\\________/\\\_______/\\\\\_______/\\\\\\\\\\\\_____/\\\________/\\\_____/\\\\\\\\\____        
 _\/\\\_______\/\\\_____/\\\///\\\____\/\\\////////\\\__\/\\\_____/\\\//____/\\\\\\\\\\\\\__       
  _\//\\\______/\\\____/\\\/__\///\\\__\/\\\______\//\\\_\/\\\__/\\\//______/\\\/////////\\\_      
   __\//\\\____/\\\____/\\\______\//\\\_\/\\\_______\/\\\_\/\\\\\\//\\\_____\/\\\_______\/\\\_     
    ___\//\\\__/\\\____\/\\\_______\/\\\_\/\\\_______\/\\\_\/\\\//_\//\\\____\/\\\\\\\\\\\\\\\_    
     ____\//\\\/\\\_____\//\\\______/\\\__\/\\\_______\/\\\_\/\\\____\//\\\___\/\\\/////////\\\_   
      _____\//\\\\\_______\///\\\__/\\\____\/\\\_______/\\\__\/\\\_____\//\\\__\/\\\_______\/\\\_  
       ______\//\\\__________\///\\\\\/_____\/\\\\\\\\\\\\/___\/\\\______\//\\\_\/\\\_______\/\\\_ 
        _______\///_____________\/////_______\////////////_____\///________\///__\///________\///__
`

var (
	appCmd *exec.Cmd
)

func main() {
	fmt.Println(Cyan + banner + Reset)
	log.Println(Blue + "Vodka Watching for Empty Glasses..." + Reset)

	buildAndRun()

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	err = watcher.Add(".")
	if err != nil {
		log.Fatal(err)
	}

	var timer *time.Timer

	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}

			if strings.HasSuffix(event.Name, ".go") {
				if timer != nil {
					timer.Stop()
				}
				timer = time.AfterFunc(500*time.Millisecond, func() {
					log.Printf(Yellow+"File changed: %s. Rebuilding...\n"+Reset, event.Name)
					buildAndRun()
				})
			}
		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			log.Println(Red+"Watcher Error: "+Reset, err)
		}
	}
}

func buildAndRun() {
	if appCmd != nil && appCmd.Process != nil {
		log.Println(Gray + "Stopping old process..." + Reset)
		appCmd.Process.Kill()
		appCmd.Wait()
	}

	binaryPath := filepath.Join("tmp", "vodka-build")
	if runtime.GOOS == "windows" {
		binaryPath += ".exe"
	}

	buildCmd := exec.Command("go", "build", "-o", binaryPath, ".")
	buildCmd.Stdout = os.Stdout
	buildCmd.Stderr = os.Stderr

	if err := buildCmd.Run(); err != nil {
		log.Println(Red + "Build failed. Waiting for changes..." + Reset)
		return
	}

	runPath := "." + string(filepath.Separator) + binaryPath
	appCmd = exec.Command(runPath)
	appCmd.Stdout = os.Stdout
	appCmd.Stderr = os.Stderr

	if err := appCmd.Start(); err != nil {
		log.Println(Red+"Failed to start app:"+Reset, err)
		return
	}

	log.Println(Green + "App is running!" + Reset)
}
