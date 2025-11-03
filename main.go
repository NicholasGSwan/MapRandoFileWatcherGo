package main

import (
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/fsnotify/fsnotify"
)

const (
	downloadsDir      string = "E:/downloads"
	mapRandosDir      string = "E:/downloads/maprandos"
	mapRandoFileRegex string = "(map-rando-)([a-zA-Z0-9]+)(.sfc)"
)

func main() {
	watcher, err := fsnotify.NewWatcher()
	regex, err2 := regexp.Compile(mapRandoFileRegex)
	if err != nil {
		log.Fatal(err)
	}
	if err2 != nil {
		log.Fatal(err2)
	}
	defer watcher.Close()

	// Start listening for events.
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				log.Println("event:", event)
				if event.Has(fsnotify.Create) && regex.MatchString(event.Name) {
					strArr := strings.Split(event.Name, "\\")
					fName := strArr[len(strArr)-1]
					log.Println("modified file:", event.Name)
					log.Println("file without path is:", fName)
					err3 := os.Rename(downloadsDir+"\\"+fName, mapRandosDir+"\\"+fName)
					if err3 != nil {
						log.Fatal(err3)
					}
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			}
		}
	}()

	err = watcher.Add(downloadsDir)
	if err != nil {
		log.Fatal(err)
	}

	// Block main goroutine forever.
	<-make(chan struct{})

}
