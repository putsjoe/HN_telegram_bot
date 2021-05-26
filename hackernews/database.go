package hackernews

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"strings"
)

// Database handles the db connection
type Database struct {
	DB *sql.DB
}

// InitDB initiates the database connection
func InitDB() *sql.DB {
	db, err := sql.Open("sqlite3", "./data.sqlite")
	if err != nil {
		log.Fatal(err)
	}
	err = db.Ping()
	if err != nil {
		panic("Cant get DB")
	}
	return db
}

func (db Database) addItem(hn hnResponse) {
	st, err := db.DB.Prepare(
		"INSERT OR IGNORE INTO posts(ID, Title, Score, Comments, URL) " +
			"VALUES (?,?,?,?,?)")
	_, err = st.Exec(hn.ID, hn.Title, hn.Score, hn.comments, hn.URL)
	if err != nil {
		log.Println(err)
	}
}

func (db Database) checkItem(item string) bool {
	var res int
	exists := db.DB.QueryRow(
		"SELECT id FROM posts WHERE id=?", strings.TrimSpace(item))
	err := exists.Scan(&res)
	if err == sql.ErrNoRows {
		return false
	}
	return true
}

func (db Database) markRead(ids []int) {
	var sids []string
	for _, v := range ids {
		sids = append(sids, strconv.Itoa(v))
	}
	sql := fmt.Sprintf(
		"UPDATE posts SET read = 1 WHERE ID IN (%s);", strings.Join(sids, ","))
	st, err := db.DB.Prepare(sql)
	_, err = st.Exec()
	if err != nil {
		log.Println(err)
	}
	// log.Println("Marked as Read - ", strings.Join(sids, ", "))
}

func (db Database) getFive() []hnResponse {
	posts := make([]hnResponse, 0)
	rows, _ := db.DB.Query("select ID, Title, Score, URL FROM " +
		"posts WHERE read = 0 LIMIT 5;")
	var itm hnResponse
	var ids []int

	for rows.Next() {
		rows.Scan(&itm.ID, &itm.Title, &itm.Score, &itm.URL)
		posts = append(posts, itm)
		ids = append(ids, itm.ID)
	}
	go db.markRead(ids)
	return posts
}

func (db Database) savePost(userID int, postID int) {
	st, err := db.DB.Prepare("INSERT INTO user(userID, postID) VALUES(?, ?)")
	_, err = st.Exec(userID, postID)
	if err != nil {
		log.Println(err)
	}
}

func (db Database) deletePost(userID int, postID int) {
	st, err := db.DB.Prepare("DELETE FROM user WHERE userID=? AND postID=?")
	_, err = st.Exec(userID, postID)
	if err != nil {
		log.Println(err)
	}
}

func (db Database) stats() (int, int) {
	var u, t int
	exists := db.DB.QueryRow("SELECT count(*) FROM posts WHERE read = 0")
	err := exists.Scan(&u)
	if err != nil {
		fmt.Println(err)
	}
	exists = db.DB.QueryRow("SELECT count(*) FROM posts")
	err = exists.Scan(&t)
	if err != nil {
		fmt.Println(err)
	}
	return u, t
}

func (db Database) getSavedPosts(userID int) []hnResponse {

	rows, err := db.DB.Query("SELECT postID FROM user WHERE userID=?", userID)
	if err != nil {
		fmt.Println(err)
		if err == sql.ErrNoRows {
			return nil
		}
	}
	var ids []int
	var id int

	for rows.Next() {
		rows.Scan(&id)
		ids = append(ids, id)
	}

	// Convert ids to strings
	var sids []string
	for _, v := range ids {
		sids = append(sids, strconv.Itoa(v))
	}
	// Prepare statement to get all this data
	sql := fmt.Sprintf(
		"SELECT ID,Title,URL FROM posts WHERE ID IN (%s) LIMIT 30;", strings.Join(sids, ","))
	rows, err = db.DB.Query(sql)
	posts := make([]hnResponse, 0)
	var itm hnResponse

	for rows.Next() {
		rows.Scan(&itm.ID, &itm.Title, &itm.URL)
		posts = append(posts, itm)
	}
	return posts
}
