package hackernews

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

const hackerURL = "https://news.ycombinator.com/?item="
const apiURL = "https://hacker-news.firebaseio.com/v0/"
const topStories = apiURL + "topstories"
const getItem = apiURL + "item/"
const tmpfile = "/tmp/hackernews-reads.txt"

type hnResponse struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	URL       string `json:"url"`
	Score     int    `json:"score"`
	ItemType  string `json:"type"`
	Kids      []int  `json:"kids"`
	comments  int
	hackerURL string
}

func itemURL(itemID string) string {
	return getItem + strings.TrimSpace(itemID) + ".json?print=pretty"
}

func latestItems() []hnResponse {

	resp, err := http.Get(topStories + ".json?print=pretty")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	scanner := bufio.NewScanner(resp.Body)
	scanner.Scan()

	text := strings.Trim(scanner.Text(), "[]")
	ids := strings.Split(text, ", ")
	items := ids[:3]

	newItems := make([]hnResponse, 0)

	for _, item := range items {
		resp, err = http.Get(itemURL(item))
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()

		var js hnResponse
		err = json.NewDecoder(resp.Body).Decode(&js)

		js.comments = len(js.Kids)
		js.hackerURL = hackerURL + strconv.Itoa(js.ID)

		if err != nil {
			fmt.Println(err)
		}
		newItems = append(newItems, js)
	}

	return newItems
}

func PrintItems() string {
	var text string
	ni := latestItems()

	for _, each := range ni {
		text = text + each.Title + "\nComments " + strconv.Itoa(each.comments) +
			" | Score " + strconv.Itoa(each.Score) + "\n" + each.hackerURL +
			"\n\n"
	}
	return text
}
