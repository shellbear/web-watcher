package main

import (
	"github.com/bwmarrin/discordgo"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"log"
	"os"
	"os/signal"
	"syscall"
)

var db *gorm.DB
var dg *discordgo.Session

func main() {
	var err error

	if dg, err = discordgo.New("Bot " + os.Getenv("DISCORD_TOKEN")); err != nil {
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
