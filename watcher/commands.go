package watcher

import (
	"fmt"
	"log"
	"net/url"
	"strings"

	"github.com/bwmarrin/discordgo"

	"github.com/shellbear/web-watcher/models"
)

// The default handler to use when a new message is sent.
func (w *Watcher) onNewMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	args := strings.Split(strings.TrimSpace(m.Content), " ")
	if m.Author.ID == s.State.User.ID || len(args) == 0 {
		return
	}

	var err error
	switch args[0] {
	case w.Prefix + "watch":
		_, err = w.watch(s, m, args)
	case w.Prefix + "unwatch":
		_, err = w.unwatch(s, m, args)
	case w.Prefix + "watchlist":
		_, err = w.watchList(s, m, args)
	}

	if err != nil {
		log.Printf("Failed to execute command '%s'. Error: %s\n", args[0], err)
	}
}

// The watch discord command handler.
// Used to add a task to the list.
func (w *Watcher) watch(s *discordgo.Session, m *discordgo.MessageCreate, args []string) (*discordgo.Message, error) {
	if len(args) != 2 {
		return s.ChannelMessageSend(m.ChannelID, m.Author.Mention()+" Usage: watch URL")
	}

	u, err := url.ParseRequestURI(args[1])
	if err != nil || u.String() == "" {
		return s.ChannelMessageSend(m.ChannelID, m.Author.Mention()+" provided URL is invalid")
	}

	var task models.Task
	if !w.DB.Where("url = ? AND guild_id = ?", u.String(), m.GuildID).First(&task).RecordNotFound() {
		return s.ChannelMessageSend(m.ChannelID, m.Author.Mention()+" this URL has already been registered")
	}

	task.URL = u.String()
	task.GuildID = m.GuildID
	task.ChannelID = m.ChannelID
	if err := w.DB.Create(&task).Error; err != nil {
		return nil, err
	}

	w.NewTask(&task)

	return s.ChannelMessageSend(m.ChannelID, m.Author.Mention()+" successfully registered URL")
}

// The unwatch discord command handler.
// Used to remove a task from the list.
func (w *Watcher) unwatch(s *discordgo.Session, m *discordgo.MessageCreate, args []string) (*discordgo.Message, error) {
	if len(args) != 2 {
		return s.ChannelMessageSend(m.ChannelID, m.Author.Mention()+" Usage: unwatch URL")
	}

	u, err := url.ParseRequestURI(args[1])
	if err != nil {
		return s.ChannelMessageSend(m.ChannelID, m.Author.Mention()+" provided URL is invalid")
	}

	var task models.Task
	if w.DB.Where("url = ? AND guild_id = ?", u.String(), m.GuildID).First(&task).RecordNotFound() {
		return s.ChannelMessageSend(m.ChannelID, m.Author.Mention()+" URL doesn't exist")
	}

	taskName := task.URL + task.ChannelID
	if cancel, ok := w.Tasks[taskName]; ok {
		cancel()
		delete(w.Tasks, taskName)
	}

	if err := w.DB.Delete(&task).Error; err != nil {
		return nil, err
	}

	return s.ChannelMessageSend(m.ChannelID, m.Author.Mention()+" successfully deleted URL")
}

// The watchlist discord command handler.
// Used to retrieve the list of tasks.
func (w *Watcher) watchList(s *discordgo.Session, m *discordgo.MessageCreate, args []string) (*discordgo.Message, error) {
	var tasks []models.Task

	if err := w.DB.Where("guild_id = ?", m.GuildID).Find(&tasks).Error; err != nil {
		return nil, err
	}

	if len(tasks) == 0 {
		return s.ChannelMessageSend(m.ChannelID, m.Author.Mention()+"There is no registered URL. Add one with `watch` command")
	}

	var urls []string
	for i, task := range tasks {
		urls = append(urls, fmt.Sprintf("%d - %s", i+1, task.URL))
	}

	return s.ChannelMessageSend(m.ChannelID, m.Author.Mention()+"\n"+strings.Join(urls, "\n"))
}

// The ready discord handler.
// Used to set the bot status.
func (w *Watcher) onReady(discord *discordgo.Session, ready *discordgo.Ready) {
	if err := discord.UpdateStatus(0, "Looking at other people's websites"); err != nil {
		log.Fatalln("Error attempting to set my status,", err)
	}

	log.Printf("Web-watcher has started on %d servers\n", len(discord.State.Guilds))
	log.Printf("Inspecting websites every %d minutes", int(w.WatchInterval.Minutes()))
}
