package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/jinzhu/gorm/dialects/sqlite"

	"github.com/shellbear/web-watcher/watcher"
)

var (
	watchInterval int64
	changeRatio   float64
	discordToken  string
	commandPrefix string
)

func init() {
	flag.Int64Var(&watchInterval, "interval", int64(time.Hour.Minutes()), "The watcher interval in minutes")
	flag.StringVar(&discordToken, "token", "", "Discord token")
	flag.Float64Var(&changeRatio, "ratio", 1.0, "Changes detection ratio")
	flag.StringVar(&commandPrefix, "prefix", "!", "The discord commands prefix")

	flag.Usage = usage
	flag.Parse()

	if discordToken == "" {
		if token := os.Getenv("DISCORD_TOKEN"); token == "" {
			log.Fatalln("you must provide a discord token")
		} else {
			discordToken = token
		}
	}

	if changeRatio <= 0 || changeRatio > 1 {
		log.Fatalln("change ratio must be between 0 and 1")
	}
}

func usage() {
	fmt.Println("Web-watcher discord Bot.\n\nOptions:")
	flag.PrintDefaults()
}

func main() {
	instance, err := watcher.New(time.Duration(watchInterval)*time.Minute, changeRatio, discordToken, commandPrefix)
	if err != nil {
		log.Fatalln(err)
	}

	defer func() {
		instance.DB.Close()
		instance.Session.Close()
	}()

	if err := instance.Run(); err != nil {
		log.Fatalln("failed to run tasks:", err)
	}

	log.Println("Bot is now running.  Press CTRL-C to exit.")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
	log.Println("Gracefully stopping bot...")
}
