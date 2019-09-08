package main

import (
	"flag"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var db *gorm.DB
var dg *discordgo.Session

var discordToken = flag.String("token", os.Getenv("DISCORD_TOKEN"), "Discord token")
var watchDelay = flag.Int64("delay", int64(time.Hour.Minutes()), "Watch delay in minutes")

func usage() {
	fmt.Fprintf(os.Stderr, "Web-watcher discord Bot.\n\nOptions:\n")
	flag.PrintDefaults()
}

func main() {
	var err error

	flag.Usage = usage

	flag.Parse()

	if *discordToken == "" {
		log.Fatalln("You must provide a discord token.")
	}

	if dg, err = discordgo.New("Bot " + *discordToken); err != nil {
		log.Fatalln(err)
	}
	if db, err = gorm.Open("sqlite3", "db.sqlite"); err != nil {
		log.Fatalln("failed to connect database")
	}

	defer db.Close()
	db.AutoMigrate(&Website{})

	if err != nil {
		log.Fatalln("Error creating Discord session,", err)
	}

	dg.AddHandler(messageCreate)
	dg.AddHandlerOnce(ready)

	if err := dg.Open(); err != nil {
		log.Fatalln("Error opening connection,", err)
	}

	var websites []Website
	db.Find(&websites)

	for _, website := range websites {
		launchTask(&website)
	}

	log.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	dg.Close()
}
