package main

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

func crawlWebsite(ctx context.Context, website *Website) error {
	log.Println("Crawling website", website.Url)

	resp, err := http.Get(website.Url)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	hasher := sha256.New()
	if _, err := hasher.Write(content); err != nil {
		return err
	}

	hash := hex.EncodeToString(hasher.Sum(nil))

	if website.Hash == "" {
		db.Model(website).Update("hash", hash)
		log.Println("Saved hash for", website.Url)
	} else if hash != website.Hash {
		_, err := dg.ChannelMessageSend(website.ChannelID,
			fmt.Sprintf("%s has been updated! Last update %x", website.Url, website.UpdatedAt))

		if err != nil {
			return err
		}

		db.Model(website).Update("hash", hash)
		log.Println("Hash differs for", website.Url)
	} else {
		log.Println("Got same hash for", website.Url)
	}

	select {
	case <-time.After(time.Hour):
		db.Find(website, website.ID)
		return crawlWebsite(ctx, website)
	case <-ctx.Done():
		log.Println("Stopped task for", website.Url)
		return nil
	}
}
