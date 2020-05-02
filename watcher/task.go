package watcher

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/google/brotli/go/cbrotli"

	"github.com/OneOfOne/xxhash"
	"golang.org/x/net/html"

	"github.com/shellbear/web-watcher/models"
)

const updateFormat = "2006-01-02 15:04:05"

// Add a new task to the list and run it.
// If a task already exits for the given URL then cancel it and override it.
func (w *Watcher) NewTask(task *models.Task) {
	taskName := task.URL + task.ChannelID
	ctx, cancel := context.WithCancel(context.Background())

	if cancel, ok := w.Tasks[taskName]; ok {
		cancel()
	}

	w.Tasks[taskName] = cancel

	go func() {
		if err := w.runTask(ctx, task); err != nil {
			log.Printf("Failed to run task for %s. Error: %s\n", task.URL, err)
		}
	}()
}

// Extract body and generate a unique hash for the web page.
// This hash will be used to speed up future page comparisons.
func (w *Watcher) getHash(resp *http.Response) (string, []byte, error) {
	xxHash := xxhash.New64()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", nil, err
	}

	// Parse page as HTML.
	doc, err := html.Parse(bytes.NewBuffer(body))
	if err != nil {
		return "", nil, err
	}

	// Try to parse only the content of body by default.
	if bn, err := getBody(doc); err == nil {
		var buf bytes.Buffer

		if err := html.Render(io.Writer(&buf), bn); err != nil {
			return "", nil, err
		}

		body = buf.Bytes()
		// Create a unique hash for the page content.
		if _, err := xxHash.Write(body); err != nil {
			return "", nil, err
		}

		return strconv.FormatUint(xxHash.Sum64(), 10), body, nil
	}

	// Create a unique hash for the page content.
	if _, err := xxHash.Write(body); err != nil {
		return "", nil, err
	}

	return strconv.FormatUint(xxHash.Sum64(), 10), body, nil
}

// Analyze page structure, extract tags and check difference ratio between changes.
func (w *Watcher) hasChanged(task *models.Task, body []byte, hash string) (bool, error) {
	// Page is identical to the previous one, we skip further checks
	if task.Hash == hash {
		return false, nil
	}

	// Check if web page has changed.
	updated, err := w.checkChanges(task, body)
	if err != nil {
		return false, err
	}

	if updated {
		// Encode the body to decrease the size in database.
		encodedBody, err := cbrotli.Encode(body, cbrotli.WriterOptions{
			Quality: 11,
		})
		if err != nil {
			return false, err
		}

		if err := w.updateTask(task, hash, encodedBody); err != nil {
			return false, err
		}

		return true, nil
	}

	return true, nil
}

func (w *Watcher) analyzeChanges(task *models.Task) error {
	resp, err := http.Get(task.URL)
	if err != nil {
		return fmt.Errorf("failed to fetch task URL: %s", err)
	}

	defer resp.Body.Close()

	hash, body, err := w.getHash(resp)
	if err != nil {
		return fmt.Errorf("failed to parse and generate hash for page. Error: %s", err)
	}

	updated, err := w.hasChanged(task, body, hash)
	if err != nil {
		return fmt.Errorf("failed to check page changes: %s", err)
	}

	if updated {
		log.Printf("%s has been updated.\n", task.URL)
	}

	return nil
}

// Run the task every X minutes.
func (w *Watcher) runTask(ctx context.Context, task *models.Task) error {
	log.Println("Crawling:", task.URL)

	if err := w.analyzeChanges(task); err != nil {
		fmt.Println(err)
	}

	select {
	case <-time.After(w.WatchInterval):
		if err := w.DB.Find(task, task.ID).Error; err != nil {
			return err
		}

		return w.runTask(ctx, task)
	case <-ctx.Done():
		log.Println("Stopped task for", task.URL)
		return nil
	}
}
