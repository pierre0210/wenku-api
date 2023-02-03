package user

import (
	"crypto/sha256"
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

func User(db *sql.DB) (sql.Result, error) {
	stmt, err := db.Prepare(`CREATE TABLE IF NOT EXISTS user_info(
		user_name TEXT PRIMARY KEY,
		password TEXT,
		content TEXT
	)`)
	result, err := stmt.Exec()
	return result, err
}

func CheckUser(db *sql.DB, userName string) (bool, error) {
	var state bool
	stmt, _ := db.Prepare(`SELECT user_name FROM user_info WHERE user_name = ?`)
	rows, err := stmt.Query(userName)
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

func CheckPassword(db *sql.DB, userName string, password string) (bool, error) {
	var state bool
	var dbHash [32]byte
	stmt, _ := db.Prepare(`SELECT password FROM user_info WHERE user_name = ?`)
	rows, err := stmt.Query(userName)
	if err != nil {
		return false, err
	}
	defer rows.Close()
	rows.Scan(&dbHash)
	passHash := sha256.Sum256([]byte(password))
	if passHash == dbHash {
		state = true
	} else {
		state = false
	}

	return state, nil
}
