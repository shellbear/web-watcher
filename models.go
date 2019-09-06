package main

import "github.com/jinzhu/gorm"

type Website struct {
	gorm.Model
	Url       string
	ChannelID string
	GuildID   string
	Hash      string
}
