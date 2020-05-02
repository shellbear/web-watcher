package watcher

import (
	"fmt"
	"net/http"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/jinzhu/gorm"
	"golang.org/x/net/context"

	"github.com/shellbear/web-watcher/models"
)

type Watcher struct {
	// The database instance.
	DB *gorm.DB

	// The discord session.
	Session *discordgo.Session

	// A custom HTTP client.
	Client *http.Client

	// The discord command Prefix. Usually '!'.
	Prefix string

	// The web page changes ratio. Must be between 0 and 1.
	// Every x minutes the watcher will fetch the website page and compares it with the previous version.
	// It will check changes and convert these changes to a ratio. If page are identical, this ratio is equals to 1.0,
	// and it will decrease for every detected change.
	ChangeRatio float64

	// The WatchInterval determines the interval at which the watcher will crawl web pages.
	WatchInterval time.Duration

	// A list of running tasks.
	Tasks map[string]context.CancelFunc
}

func New(watchInterval time.Duration, changeRatio float64, discordToken string, prefix string) (*Watcher, error) {
	db, err := models.New()
	if err != nil {
		return nil, fmt.Errorf("failed to init database: %s", err)
	}

	session, err := discordgo.New("Bot " + discordToken)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to discord: %s", err)
	}

	session.Client = &http.Client{Timeout: time.Second * 15}
	w := Watcher{
		DB:            db,
		Session:       session,
		Prefix:        prefix,
		WatchInterval: watchInterval,
		ChangeRatio:   changeRatio,
		Client:        session.Client,
		Tasks:         map[string]context.CancelFunc{},
	}

	session.AddHandlerOnce(w.onReady)
	session.AddHandler(w.onNewMessage)

	return &w, nil
}

// Fetch existing tasks in database and run them, then connect to Discord API.
func (w *Watcher) Run() error {
	var tasks []models.Task
	if err := w.DB.Find(&tasks).Error; err != nil {
		return fmt.Errorf("failed to fetch existing tasks: %s", err)
	}

	for i := range tasks {
		w.NewTask(&tasks[i])
	}

	return w.Session.Open()
}

// Update task in database and alert sender on Discord.
func (w *Watcher) updateTask(task *models.Task, hash string, body []byte) error {
	if _, err := w.Session.ChannelMessageSend(
		task.ChannelID,
		fmt.Sprintf("%s has been updated! Last update : %s", task.URL, task.UpdatedAt.Format(updateFormat)),
	); err != nil {
		return err
	}

	return w.DB.Model(task).Updates(&models.Task{
		Hash: hash,
		Body: body,
	}).Error
}
