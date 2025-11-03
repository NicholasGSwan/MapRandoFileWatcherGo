package main

import (
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/namsral/flag"

	"github.com/fsnotify/fsnotify"
)

const (
	mapRandoFileRegex string = "(map-rando-)([a-zA-Z0-9]+)(.sfc)"
)

func main() {
	conf := &config{}
	run(conf)

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

	//initialize config from commandline args
	conf.init(os.Args)
	//create maprando directory if not exists
	err = os.MkdirAll(conf.mapRandosDir, os.ModePerm)
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
					err3 := os.Rename(conf.downloadsDir+"\\"+fName, conf.mapRandosDir+"\\"+fName)
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

	err = watcher.Add(conf.downloadsDir)
	if err != nil {
		log.Fatal(err)
	}
	// Block main goroutine forever.
	<-make(chan struct{})
	return nil
}

type config struct {
	downloadsDir string
	mapRandosDir string
}

func (c *config) init(args []string) error {
	flags := flag.NewFlagSet(args[0], flag.ExitOnError)
	flags.String(flag.DefaultConfigFlagname, "./config.conf", "Path to config file")

	var (
		downloadsDir = flags.String("downloadsDir", "", "The downloads folder to watch for new map rando roms")
		mapRandosDir = flags.String("mapRandosDir", "", "The folder to move the map rando roms to")
	)

	if err := flags.Parse(args[1:]); err != nil {
		return err
	}

	c.downloadsDir = *downloadsDir
	c.mapRandosDir = *mapRandosDir

	return nil
}
