package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"log"
	"net/url"
	"strings"
	"time"
)

type commandProto func(*discordgo.Session, *discordgo.MessageCreate, []string) (*discordgo.Message, error)

var commands = map[string]commandProto{
	"!watch":     watch,
	"!unwatch":   unwatch,
	"!watchlist": watchList,
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

	db.Find(&websites, Website{GuildID: m.GuildID})

	if len(websites) == 0 {
		return s.ChannelMessageSend(m.ChannelID, m.Author.Mention()+"There is no registered URL. Add one with `watch` command")
	}

	var urls []string

	for i, website := range websites {
		urls = append(urls, fmt.Sprintf("%d - %s", i+1, website.Url))
	}

	return s.ChannelMessageSend(m.ChannelID, m.Author.Mention()+"\n"+strings.Join(urls, "\n"))
}

func ready(discord *discordgo.Session, ready *discordgo.Ready) {
	if err := discord.UpdateStatus(0, "Looking at other people's website"); err != nil {
		log.Fatalln("Error attempting to set my status,", err)
	}

	servers := discord.State.Guilds
	log.Printf("GoBot has started on %d servers\n", len(servers))
	log.Printf("Inspecting websites every %f minutes", (time.Duration(*watchDelay) * time.Minute).Minutes())
}
