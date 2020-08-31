package hackernews

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const hackerURL = "https://news.ycombinator.com/?item="
const apiURL = "https://hacker-news.firebaseio.com/v0/"
const topStories = apiURL + "topstories"
const getItem = apiURL + "item/"
const tmpfile = "/tmp/hackernews-reads.txt"

type HNresponse struct {
	ID       int    `json:"id"`
	Title    string `json:"title"`
	Score    int    `json:"score"`
	ItemType string `json:"type"`
	Kids     []int  `json:"kids"`
	comments int
	URL      string
}

func itemURL(itemID string) string {
	return getItem + strings.TrimSpace(itemID) + ".json?print=pretty"
}

// TextItems returns the latest items formatted for the Telegram reply
// along with the post ids so they can be marked as read
func TextItems() string {
	var text string
	db := Database{DB: InitDB()}
	ni := db.getFive()

	for _, each := range ni {
		text = text + each.Title + "\nScore " + strconv.Itoa(each.Score) +
			" | /add_" + strconv.Itoa(each.ID) + "\n" + each.URL +
			"\n\n"
	}

	return text
}

// UpdatePosts updates the database with the latest posts
func UpdatePosts() {
	getLatest()
	ticker := time.NewTicker(30 * time.Minute)
	for range ticker.C {
		getLatest()
	}
}

func getLatest() {
	db := Database{DB: InitDB()}
	resp, err := http.Get(topStories + ".json?print=pretty")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	scanner := bufio.NewScanner(resp.Body)
	scanner.Scan()

	text := strings.Trim(scanner.Text(), "[]")
	items := strings.Split(text, ", ")

	for _, item := range items {
		if db.checkItem(item) {
			continue
		}
		resp, err = http.Get(itemURL(item))
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()

		var js HNresponse
		err = json.NewDecoder(resp.Body).Decode(&js)

		js.comments = len(js.Kids)
		js.URL = hackerURL + strconv.Itoa(js.ID)

		if err != nil {
			fmt.Println(err)
		}
		db.addItem(js)
	}
	log.Println("Finished DB update")

}

// Function to keep grabbing the latest posts and adding them to a database.
// -- Filter out anything that isnt a story

// Function to return 3 - 5 unread posts and mark as read.
