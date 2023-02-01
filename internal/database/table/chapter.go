package table

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func Chapter(db *sql.DB) (sql.Result, error) {
	stmt, err := db.Prepare(`CREATE TABLE IF NOT EXISTS chapter_info(
		chapter_id TEXT PRIMARY KEY,
		title TEXT,
		content TEXT
	)`)
	result, err := stmt.Exec()
	return result, err
}

func ChapterExists(db *sql.DB, aid int, vid int, cid int) (bool, error) {
	var state bool
	stmt, _ := db.Prepare(`SELECT chapter_id FROM chapter_info WHERE chapter_id=?`)
	rows, err := stmt.Query(fmt.Sprintf("%d-%d-%d", aid, vid, cid))
	log.Println(rows)
	if err != nil {
		return false, err
	}
	if !rows.Next() {
		state = false
	} else {
		state = true
	}
	return state, nil
}
