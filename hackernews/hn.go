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

const hackerURL = "https://news.ycombinator.com/item?id="
const apiURL = "https://hacker-news.firebaseio.com/v0/"
const topStories = apiURL + "topstories"
const getItem = apiURL + "item/"
const tmpfile = "/tmp/hackernews-reads.txt"

type hnResponse struct {
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

		var js hnResponse
		err = json.NewDecoder(resp.Body).Decode(&js)

		if js.Score < 11 {
			continue
		}
		js.comments = len(js.Kids)
		js.URL = hackerURL + strconv.Itoa(js.ID)

		if err != nil {
			fmt.Println(err)
		}
		db.addItem(js)
	}
	log.Println("Finished DB update")

}

// UpdatePosts updates the database with the latest posts
func UpdatePosts() {
	getLatest()
	ticker := time.NewTicker(1 * time.Hour)
	for range ticker.C {
		getLatest()
	}
}

// SavePost saves the given post ID for the user to read later
func SavePost(userID int, postID int) {
	db := Database{DB: InitDB()}
	db.savePost(userID, postID)
}

// DeletePost removes the given id from the user's list
func DeletePost(userID int, postID int) {
	db := Database{DB: InitDB()}
	db.deletePost(userID, postID)
}

// UnreadItems returns the number of items in the database unread and total
func UnreadItems() (int, int) {
	db := Database{DB: InitDB()}
	return db.stats()
}

// TextItems returns the latest items formatted for the Telegram reply
// along with the post ids so they can be marked as read
func TextItems() string {
	var text string
	db := Database{DB: InitDB()}
	ni := db.getFive()

	for _, each := range ni {
		text = text + each.Title + "\nScore " + strconv.Itoa(each.Score) +
			" | /add_" + strconv.Itoa(each.ID) + "\n\n"
	}

	return text
}

// GetSavedPosts returns the user's saved posts
func GetSavedPosts(userID int) string {
	db := Database{DB: InitDB()}
	posts := db.getSavedPosts(userID)

	if len(posts) == 0 {
		return "No Saved Posts"
	}

	var text string
	for _, post := range posts {
		text = text + post.Title + " | /del_" + strconv.Itoa(post.ID) +
			"\n" + "<a href=\"" + post.URL + "\">HN Link</a>\n\n"
	}
	return text
}
