package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/fsnotify/fsnotify"
)

const (
	mapRandoFileRegex string = "(map-rando-)([a-zA-Z0-9]+)(.sfc)"
)

func main() {
	conf, err := loadConfig(".\\config.conf")
	if err != nil {
		os.Exit(1)
	}
	run(&conf)

}

func run(conf *config) error {
	watcher, err := fsnotify.NewWatcher()
	regex, err2 := regexp.Compile(mapRandoFileRegex)

	if err != nil {
		log.Fatal(err)
	}
	if err2 != nil {
		log.Fatal(err2)
	}
	defer watcher.Close()
	print(conf.LogDir, " is the log dir in the config\n")
	print(conf.MapRandosDir, " is the maprando dir in the config\n")
	// err = os.MkdirAll(conf.LogDir, os.ModePerm)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	fmt.Println(conf.LogDir+"\\mapRando.log", " is the current log directory")
	logFile, err := os.Create(conf.LogDir + "\\mapRando.log")
	// logFile, err := os.OpenFile(conf.logDir+"\\mapRando.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		log.Fatal(err)
	}
	defer logFile.Close()
	log.SetOutput(logFile)
	//initialize config from commandline args

	//create maprando directory if not exists
	err = os.MkdirAll(conf.MapRandosDir, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}

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
					err3 := os.Rename(conf.DownloadsDir+"\\"+fName, conf.MapRandosDir+"\\"+fName)
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

	err = watcher.Add(conf.DownloadsDir)
	if err != nil {
		log.Fatal(err)
	}
	// Block main goroutine forever.
	<-make(chan struct{})
	return nil
}

type config struct {
	DownloadsDir string `json:"downloadsDir"`
	MapRandosDir string `json:"mapRandosDir"`
	LogDir       string `json:"logDir"`
}

func loadConfig(path string) (config, error) {
	var cfg config
	data, err := os.ReadFile(path)
	if err != nil {
		return cfg, err
	}

	err = json.Unmarshal(data, &cfg)

	return cfg, err
}
