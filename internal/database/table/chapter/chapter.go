package chapter

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

func Chapter(db *sql.DB) (sql.Result, error) {
	stmt, err := db.Prepare(`CREATE TABLE IF NOT EXISTS chapter_info(
		chapter_id VARCHAR(15) PRIMARY KEY,
		title TEXT,
		content TEXT
	)`)
	result, err := stmt.Exec()
	return result, err
}

func CheckChapter(db *sql.DB, aid int, vid int, cid int) (bool, error) {
	var state bool
	stmt, _ := db.Prepare(`SELECT chapter_id FROM chapter_info WHERE chapter_id = ?`)
	rows, err := stmt.Query(fmt.Sprintf("%d-%d-%d", aid, vid, cid))
	if err != nil {
		return false, err
	}
	defer rows.Close()
	if !rows.Next() {
		state = false
	} else {
		state = true
	}

	return state, nil
}

func AddChapter(db *sql.DB, aid int, vid int, cid int, title string, content string) (sql.Result, error) {
	stmt, _ := db.Prepare(`INSERT INTO chapter_info(chapter_id, title, content) VALUES(?, ?, ?)`)
	result, err := stmt.Exec(fmt.Sprintf("%d-%d-%d", aid, vid, cid), title, content)

	return result, err
}

func GetChapter(db *sql.DB, aid int, vid int, cid int) (string, string, string, error) {
	var chapterId string
	var title string
	var content string
	stmt, _ := db.Prepare(`SELECT * FROM chapter_info WHERE chapter_id = ?`)
	rows, err := stmt.Query(fmt.Sprintf("%d-%d-%d", aid, vid, cid))
	if err != nil {
		return "", "", "", err
	}
	defer rows.Close()
	if rows.Next() {
		err = rows.Scan(&chapterId, &title, &content)
		if err != nil {
			return "", "", "", err
		}
	}
	return chapterId, title, content, nil
}
