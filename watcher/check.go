package watcher

import (
	"bytes"
	"log"

	"github.com/google/brotli/go/cbrotli"
	"github.com/pmezard/go-difflib/difflib"
	"golang.org/x/net/html"

	"github.com/shellbear/web-watcher/models"
)

func (w *Watcher) checkChanges(task *models.Task, body []byte) (bool, error) {
	if task.Body == nil {
		return true, nil
	}

	previousBody, err := cbrotli.Decode(task.Body)
	if err != nil {
		return false, err
	}

	previousHTML, err := html.Parse(bytes.NewBuffer(previousBody))
	if err != nil {
		return false, err
	}

	newHTML, err := html.Parse(bytes.NewBuffer(body))
	if err != nil {
		return false, err
	}

	matcher := difflib.NewMatcher(extractTags(previousHTML), extractTags(newHTML))
	ratio := matcher.Ratio()

	if ratio < w.ChangeRatio {
		log.Printf("Changes detected for: %s. Changes ratio: %f < %f\n", task.URL, ratio, w.ChangeRatio)
		return true, nil
	}

	log.Println("No changed detected for:", task.URL)
	return false, nil
}
