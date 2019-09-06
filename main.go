package main

import (
	"context"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"log"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

type Website struct {
	gorm.Model
	Url       string
	ChannelID string
	GuildID   string
	Hash      string
}

type commandProto func(*discordgo.Session, *discordgo.MessageCreate, []string) (*discordgo.Message, error)

var commands = map[string]commandProto{
	"!watch":     watch,
	"!unwatch":   unwatch,
	"!watchlist": watchList,
}

var tasks = map[string]context.CancelFunc{}

var db *gorm.DB
var dg *discordgo.Session

func launchTask(website *Website) {
	taskName := website.Url + website.ChannelID
	ctx, cancel := context.WithCancel(context.Background())

	val, ok := tasks[taskName]
	if ok {
		val()
	}

	tasks[taskName] = cancel
	go crawlWebsite(ctx, website)
}

func watch(s *discordgo.Session, m *discordgo.MessageCreate, args []string) (*discordgo.Message, error) {
	if len(args) != 2 {
		return s.ChannelMessageSend(m.ChannelID, m.Author.Mention()+" Usage: watch URL")
	}

	url, err := url.Parse(args[1])
	if err != nil {
		return s.ChannelMessageSend(m.ChannelID, m.Author.Mention()+" provided URL is invalid")
	}

	website := Website{
		Url:     url.String(),
		GuildID: m.GuildID,
	}

	if db.First(&website, website).RecordNotFound() {
		website.ChannelID = m.ChannelID

		db.Create(&website)
		launchTask(&website)

		return s.ChannelMessageSend(m.ChannelID, m.Author.Mention()+" successfully registered URL")
	}

	return s.ChannelMessageSend(m.ChannelID, m.Author.Mention()+" this URL has already been registered")
}

func unwatch(s *discordgo.Session, m *discordgo.MessageCreate, args []string) (*discordgo.Message, error) {
	if len(args) != 2 {
		return s.ChannelMessageSend(m.ChannelID, m.Author.Mention()+" Usage: unwatch URL")
	}

	url, err := url.Parse(args[1])
	if err != nil {
		return s.ChannelMessageSend(m.ChannelID, m.Author.Mention()+" provided URL is invalid")
	}

	var website Website

	if !db.First(&website, Website{Url: url.String(), GuildID: m.GuildID}).RecordNotFound() {
		taskName := website.Url + website.ChannelID
		if val, ok := tasks[taskName]; ok {
			val()
			delete(tasks, taskName)
		}

		db.Delete(&website)
		return s.ChannelMessageSend(m.ChannelID, m.Author.Mention()+" successfully deleted URL")
	}

	return s.ChannelMessageSend(m.ChannelID, m.Author.Mention()+" URL doesn't exist")
}

func watchList(s *discordgo.Session, m *discordgo.MessageCreate, args []string) (*discordgo.Message, error) {
	var websites []Website

	db.Find(&websites)

	if len(websites) == 0 {
		return s.ChannelMessageSend(m.ChannelID, m.Author.Mention()+"There is no registered URL. Add one with `watch` command")
	}

	var urls []string

	for i, website := range websites {
		urls = append(urls, fmt.Sprintf("%d - %s", i+1, website.Url))
	}

	return s.ChannelMessageSend(m.ChannelID, m.Author.Mention()+"\n"+strings.Join(urls, "\n"))
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	args := strings.Split(strings.TrimSpace(m.Content), " ")

	if m.Author.ID == s.State.User.ID || len(args) == 0 {
		return
	}

	if val, ok := commands[args[0]]; ok {
		if _, err := val(s, m, args); err != nil {
			log.Fatalln(err)
		}
	}
}

func ready(discord *discordgo.Session, ready *discordgo.Ready) {
	if err := discord.UpdateStatus(0, "Looking at other people's website"); err != nil {
		log.Fatalln("Error attempting to set my status,", err)
	}
	servers := discord.State.Guilds
	log.Printf("GoBot has started on %d servers\n", len(servers))
}

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
