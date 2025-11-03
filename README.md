# MapRandoFileWatcherGo
Super Metroid Map Randomizer file watcher in Golang

## About this project
When I started playing map randos, I realized I didn't like scrolling through my default downloads folder to get down to the rom I was looking for.  I figured I could make a little program to watch for the map rando files and move them to another subdirectory.  Originally, I made one in Java and set it to run on startup.  While that one works perfectly fine, I've been wanting to play with Golang and figured this would be a good starting point.

## How to use
Compile the code into an executable, then create a text file named `config.conf` in the same directory as the executable with values like the ones seen in the `config.conf.example` file.  Be sure to set the directories as they are on YOUR machine.  Run the executable, and then start downloading your map rando roms! 