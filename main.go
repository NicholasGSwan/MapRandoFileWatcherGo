package main

import (
	"io"
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

	// println("Current maprandos dir is ", conf.downloadsDir)
	// println("Current log file dir", conf.logFile)

	logFile, err := os.Create(conf.logFile + "\\mapRando.log")
	if err != nil {
		log.Fatal(err)
	}
	defer logFile.Close()
	if conf.externalLoggingOn {
		log.SetOutput(logFile)
	}

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
				if event.Has(fsnotify.Write) && regex.MatchString(event.Name) {
					strArr := strings.Split(event.Name, "\\")
					fName := strArr[len(strArr)-1]
					log.Println("modified file:", event.Name)
					log.Println("file without path is:", fName)

					//open destination file
					destFile, err := os.Create(conf.mapRandosDir + "\\" + fName)
					if err != nil {
						print("Could not open destination file: ", conf.mapRandosDir+"\\"+fName)
						log.Fatal(err)
					}
					defer destFile.Close()

					//open source file
					srcFile, err := os.Open(conf.downloadsDir + "\\" + fName)
					if err != nil {
						print("Could not open source file: ", conf.downloadsDir+"\\"+fName)
					}
					defer srcFile.Close()
					srcFile.Seek(0, io.SeekStart)
					//copy file from source to destination
					//log.Printf("Copying src file %s to dest file %s", srcFile.Name(), destFile.Name())
					bytesCopied, err3 := io.Copy(destFile, srcFile)
					for i := 0; i < 5 && bytesCopied == 0; i++ {
						bytesCopied, err3 = io.Copy(destFile, srcFile)
					}

					if err3 != nil {
						log.Fatal(err3)
					}

					err = destFile.Sync()
					if err != nil {
						log.Print("failed to sync destination file ", err)
						log.Fatal(err)
					}

					srcFile.Close()
					if err = os.Remove(srcFile.Name()); err != nil {
						log.Fatal("file copied, but failed to remove original: %w", err)
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
	downloadsDir      string
	mapRandosDir      string
	logFile           string
	externalLoggingOn bool
}

func (c *config) init(args []string) error {
	flags := flag.NewFlagSet(args[0], flag.ExitOnError)
	flags.String(flag.DefaultConfigFlagname, "./config.conf", "Path to config file")

	var (
		downloadsDir      = flags.String("downloadsDir", "", "The downloads folder to watch for new map rando roms")
		mapRandosDir      = flags.String("mapRandosDir", "", "The folder to move the map rando roms to")
		logFile           = flags.String("logFile", "./", "where logFile will be written")
		externalLoggingOn = flags.Bool("externalLoggingOn", false, "Log to an external log file")
	)

	if err := flags.Parse(args[1:]); err != nil {
		return err
	}

	c.downloadsDir = *downloadsDir
	c.mapRandosDir = *mapRandosDir
	c.logFile = *logFile
	c.externalLoggingOn = *externalLoggingOn

	return nil
}
