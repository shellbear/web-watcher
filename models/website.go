package models

import (
	"github.com/jinzhu/gorm"
)

// The default model used to store data about a web page.
type Task struct {
	*gorm.Model

	URL       string
	ChannelID string
	GuildID   string
	Hash      string
	Body      []byte
}
