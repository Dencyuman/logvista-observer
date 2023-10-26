package main

import (
	"fmt"
	"github.com/Dencyuman/logvista-observer/internal"
	"github.com/fsnotify/fsnotify"
	"log"
	"os"
)

var version = "dev"

func main() {
	fmt.Printf("=== Logvista Observer v-%s ===\n", version)
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		return
	}
	logvistaDirPath, err := internal.CreateLogvistaDir(homeDir)

	// Create new watcher.
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	// Start listening for events.
	fmt.Printf("監視を開始します: %s\n", logvistaDirPath)

	go internal.WatchFiles(watcher, logvistaDirPath)

	// Add a path.
	err = watcher.Add(logvistaDirPath)
	if err != nil {
		log.Fatal(err)
	}

	// Block main goroutine forever.
	<-make(chan struct{})
}
